package handler

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"shopa/db"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lithammer/shortuuid/v3"
	"github.com/vincent-petithory/dataurl"
)

type InputInsertProduct struct {
	Product        db.Product                `json:"product"`
	InventoryLines []db.ProductInventoryLine `json:"inventory_lines"`
	PictureLinks   []db.ProductPicture       `json:"picture_links"`
	FAQS           []db.FAQ                  `json:"faqs"`
}

func DeleteProduct(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	authAvatar := c.MustGet("authAvatar").(string)
	productID := c.Param("productID")

	if productID == "" {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "productID not present")
		return
	}

	productOwner, err := db.ProductOwner(c.Request.Context(), dbConn, productID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "determining details owner: "+err.Error())
		return
	}

	if productOwner != authAvatar {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 403, "product owner mismatch")
		return
	}

	if err := db.DeleteProduct(c.Request.Context(), dbConn, productID); err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "deleting product: "+err.Error())
		return
	}

	c.Status(http.StatusOK)
	log.WithFields(log.Fields{
		"IP":      c.ClientIP(),
		"product": productID,
		"owner":   authAvatar,
	}).Info("deleted product")
	c.JSON(200, gin.H{
		"message": "deleted product",
	})
}

func InsertProduct(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	authAvatar := c.MustGet("authAvatar").(string)
	productDetails := InputInsertProduct{}
	err := c.ShouldBindJSON(&productDetails)
	if err != nil {
		if fe, ok := err.(validator.ValidationErrors); ok {
			err1 := fe[0]
			logRespondError(c, log.Fields{
				"IP": c.ClientIP(),
			}, 500, translate(err1), gin.H{
				"field": err1.Field(),
			})
			return
		}
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding product:"+err.Error())
		return
	}

	if len(productDetails.InventoryLines) == 0 {
		if err != nil {
			logRespondError(c, log.Fields{
				"IP": c.ClientIP(),
			}, 500, "must have items")
			return
		}
	}

	if productDetails.Product.DiscountActive && productDetails.Product.DiscountedPrice < 100 {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "Discount price cannot be below 100.")
		return
	}

	hasNoCopyableItem := false
	for _, line := range productDetails.InventoryLines {
		if !line.Copyable {
			hasNoCopyableItem = true
			break
		}
	}

	if !hasNoCopyableItem {
		var hasDemoItem bool
		for _, line := range productDetails.InventoryLines {
			if line.DemoItem == true {
				hasDemoItem = true
				break
			}
		}

		if !hasDemoItem {
			logRespondError(c, log.Fields{
				"IP": c.ClientIP(),
			}, 500, "must have demo item")
			return
		}
	}

	var hasNormalItem bool
	for _, line := range productDetails.InventoryLines {
		if line.DemoItem == false {
			hasNormalItem = true
			break
		}
	}

	if !hasNormalItem {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "must have a full-version item")
		return
	}

	if len(productDetails.PictureLinks) == 0 {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "must have a product picture")
		return
	}

	productOwner := ""
	if productDetails.Product.ID != "" {
		productOwner, err = db.ProductOwner(c.Request.Context(), dbConn, productDetails.Product.ID)
		if err != nil {
			logRespondError(c, log.Fields{
				"IP": c.ClientIP(),
			}, 500, "finding product owner:"+err.Error())
			return
		}
	}

	if productOwner != authAvatar && productDetails.Product.ID != "" {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "cannot change product you do not own")
		return
	}

	var inventoryItemIDs []string
	for _, i := range productDetails.InventoryLines {
		inventoryItemIDs = append(inventoryItemIDs, i.InventoryItemID)
	}

	nonCopyableInventoryItemsInAnotherProduct, err := db.NonCopyableInventoryItemsInAnotherProduct(c.Request.Context(), dbConn, inventoryItemIDs, authAvatar, productDetails.Product.ID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "finding no-copy items:"+err.Error())
		return
	}

	if len(nonCopyableInventoryItemsInAnotherProduct) > 0 {
		var names []string
		for _, item := range nonCopyableInventoryItemsInAnotherProduct {
			names = append(names, item.Name)
		}
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "cannot list product that contains non-copy items in another product: "+strings.Join(names, ","))
		return
	}

	productID, err := db.UpsertProduct(c.Request.Context(), dbConn, productDetails.Product, productDetails.InventoryLines, productDetails.PictureLinks, productDetails.FAQS)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "upserting product:"+err.Error())
		return
	}

	productDetails.Product.ID = productID

	log.WithFields(log.Fields{
		"IP":             c.ClientIP(),
		"productDetails": productDetails,
	}).Info("inserted/updated product")
	c.JSON(200, gin.H{
		"message":    "inserted/updated product",
		"product_id": productID,
	})
}

func FrontpageProductPreviews(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	pps, err := db.FrontpageProductPreviews(c.Request.Context(), dbConn)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "getting product previews:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP": c.ClientIP(),
	}).Info("got product previews")
	c.JSON(200, pps)
}

func DiscountedProductsFrontpages(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	pps, err := db.DiscountedProductsFrontpages(c.Request.Context(), dbConn)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "getting discounted product previews:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP": c.ClientIP(),
	}).Info("got discounted product previews")
	c.JSON(200, pps)
}

func GetProductInventoryItems(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	info := struct {
		ProductID string `json:"product_id"`
	}{}
	err := c.ShouldBindJSON(&info)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding info:"+err.Error())
		return
	}

	products, err := db.GetProductInventoryItems(c.Request.Context(), dbConn, info.ProductID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "getting avatar product from db:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":      c.ClientIP(),
		"product": products,
	}).Info("got avatar avatar product inventory items")
	c.JSON(200, products)
}

func GetAvatarProduct(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)

	authAvatar := ""
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && strings.Index(authHeader, ".") != -1 {
		split := strings.Split(authHeader, ".")
		//hash := split[0]
		authAvatar = split[1]
	}

	info := struct {
		ProductID string `json:"product_id"`
	}{}
	err := c.ShouldBindJSON(&info)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding info:"+err.Error())
		return
	}

	product, err := db.GetAvatarProduct(c.Request.Context(), dbConn, info.ProductID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "getting avatar product from db:"+err.Error())
		return
	}

	pictures, err := db.ProductPictures(c.Request.Context(), dbConn, info.ProductID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "getting avatar product pictures from db:"+err.Error())
		return
	}

	picturesLoves := []db.ProductPictureLove{}
	picturesLoves, err = db.ProductPicturesLoves(c.Request.Context(), dbConn, pictures, authAvatar)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "getting avatar product pictures loves from db:"+err.Error())
		return
	}

	type PicturesWithLoves struct {
		ID            string `json:"id"`
		ProductID     string `json:"product_id"`
		Link          string `json:"link"`
		Index         int    `json:"index"`
		LovedByAvatar bool   `json:"loved_by_avatar"`
		Loves         int    `json:"loves"`
		PictureID     string `json:"picture_id"`
	}

	var picturesWithLoves []PicturesWithLoves

	for _, picture := range pictures {
		var loveData db.ProductPictureLove
		for _, love := range picturesLoves {
			if picture.ID == love.PictureID {
				loveData = love
				break
			}
		}
		picturesWithLoves = append(picturesWithLoves, PicturesWithLoves{
			ID:            picture.ID,
			ProductID:     picture.ProductID,
			Link:          picture.Link,
			Index:         picture.Index,
			LovedByAvatar: loveData.LovedByAvatar,
			Loves:         loveData.Loves,
			PictureID:     loveData.PictureID,
		})
	}

	bought, err := db.BoughtProduct(c.Request.Context(), dbConn, info.ProductID, authAvatar)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "getting boughtness:"+err.Error())
		return
	}
	canWriteReview := bought && product.Owner != authAvatar

	blogsByAuthAvatar, err := db.BlogsByAvatarKey(c.Request.Context(), dbConn, authAvatar)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP":   c.ClientIP(),
			"blog": info,
		}, 400, "determining if blogger has already blogged in product:"+err.Error())
		return
	}
	for _, blog := range blogsByAuthAvatar {
		if blog.ProductID == info.ProductID {
			canWriteReview = false
		}
	}

	prePosts, err := db.ProductBlogPosts(c.Request.Context(), dbConn, info.ProductID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 400, "getting blog prePosts:"+err.Error())
		return
	}

	postIDsLovedByAvatar, err := db.BlogPostIDsLovedByAvatar(c.Request.Context(), dbConn, authAvatar)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 400, "getting blogs loved by avatar:"+err.Error())
		return
	}

	// convert postIDsLovedByAvatar to a map
	postIDsLovedByAvatarMap := map[string]struct{}{}
	for _, s := range postIDsLovedByAvatar {
		postIDsLovedByAvatarMap[s] = struct{}{}
	}

	type PostWithLove struct {
		ID              string `json:"id"`
		Owner           string `json:"owner"`
		ProductID       string `json:"product_id"`
		ContentLink     string `json:"content_link"`
		PictureLink     string `json:"picture_link"`
		Type            int    `json:"type"`
		OwnerLegacyName string `json:"owner_legacy_name"`
		LoveCount       int    `json:"love_count"`
		Loved           bool   `json:"loved"`
	}
	posts := make([]PostWithLove, 0, len(prePosts))
	for _, post := range prePosts {
		if _, ok := postIDsLovedByAvatarMap[post.ID]; ok {
			posts = append(posts, PostWithLove{
				ID:              post.ID,
				Owner:           post.Owner,
				ProductID:       post.ProductID,
				ContentLink:     post.ContentLink,
				PictureLink:     post.PictureLink,
				Type:            post.Type,
				OwnerLegacyName: post.OwnerLegacyName,
				LoveCount:       post.LoveCount,
				Loved:           true,
			})
			continue
		}
		posts = append(posts, PostWithLove{
			ID:              post.ID,
			Owner:           post.Owner,
			ProductID:       post.ProductID,
			ContentLink:     post.ContentLink,
			PictureLink:     post.PictureLink,
			Type:            post.Type,
			LoveCount:       post.LoveCount,
			OwnerLegacyName: post.OwnerLegacyName,
			Loved:           false,
		})
	}

	questions, err := db.ProductsQuestionsAndAnswers(c.Request.Context(), dbConn, info.ProductID, authAvatar)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 400, "getting questions and answers:"+err.Error())
		return
	}

	faqs, err := db.GetProductFAQs(c.Request.Context(), dbConn, info.ProductID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 400, "getting questions and answers:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":      c.ClientIP(),
		"product": product,
	}).Info("got avatar product")
	c.JSON(200, gin.H{
		"product":          product,
		"pictures":         picturesWithLoves,
		"can_write_review": canWriteReview,
		"bought":           bought,
		"blog_posts":       posts,
		"questions":        questions,
		"faqs":             faqs,
	})
}

func GetAvatarProducts(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	userInfo := struct {
		AvatarKey string `json:"avatar_key"`
	}{}
	err := c.ShouldBindJSON(&userInfo)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding userInfo:"+err.Error())
		return
	}

	products, err := db.GetAvatarProducts(c.Request.Context(), dbConn, userInfo.AvatarKey)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "getting avatar products from db:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":       c.ClientIP(),
		"products": products,
	}).Info("got avatar products")
	c.JSON(200, products)
}

type File struct {
	Name  string
	Bytes []byte
}

// UploadPicture accepts a jpeg, jpg, or png, uploads it, then returns the link to it.
func UploadPicture(c *gin.Context) {
	s3Client := c.MustGet("s3").(*s3.S3)
	pictureType := c.Param("pictureType")
	picProps := struct {
		maxWidth        int
		maxHeight       int
		minWidth        int
		minHeight       int
		squareDimension bool
	}{}

	if pictureType == "turntable" {
		picProps.maxWidth = 10000
		picProps.maxHeight = 1200
		picProps.minWidth = 100
		picProps.minHeight = 100
	}

	if pictureType == "product" {
		picProps.maxWidth = 500
		picProps.maxHeight = 500
		picProps.minWidth = 100
		picProps.minHeight = 100
	}

	if pictureType == "banner" {
		picProps.maxWidth = 1120
		picProps.maxHeight = 280
		picProps.minWidth = 320
		picProps.minHeight = 80
	}

	preFile, err := c.MultipartForm()
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "doing a formfile:"+err.Error())
		return
	}

	formFile, ok := preFile.Value["file"]
	if !ok {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "did not find file in form")
		return
	}

	if len(formFile) < 1 {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "did not find file in form")
		return
	}
	file := formFile[0]

	indexData := strings.Index(file, "data:")
	indexSemicolon := strings.Index(file, ";")
	mimeType := file[indexData+5 : indexSemicolon]
	if mimeType != "image/jpeg" && mimeType != "image/png" && mimeType != "image/jpg" {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "incorrect file type")
		return
	}

	fileExt := strings.Split(mimeType, "/")[1]
	fileName := shortuuid.New() + "." + fileExt

	dataURL, err := dataurl.DecodeString(file)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "failure parsing datauri to image: "+err.Error())
		return
	}

	object := s3.PutObjectInput{
		Bucket:   aws.String("shopa"),
		Key:      aws.String(fileName),
		Body:     bytes.NewReader(dataURL.Data),
		ACL:      aws.String("public-read"),
		Metadata: map[string]*string{},
	}
	_, err = s3Client.PutObject(&object)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "failure uploading to cloud: "+err.Error())
		return
	}
	link := "https://shopa.nyc3.digitaloceanspaces.com/" + fileName

	log.WithFields(log.Fields{
		"IP":    c.ClientIP(),
		"links": []string{link},
	}).Info("uploaded picture")
	c.JSON(200, gin.H{
		"message": "inserted/updated product",
		"links":   []string{link},
	})
}

func headersToFile(header *multipart.FileHeader, squareDimRequired bool, maxWidth, maxHeight, minWidth, minHeight int) (File, error) {
	file, err := header.Open()
	defer file.Close()
	if err != nil {
		return File{}, err
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return File{}, err
	}

	filetype := http.DetectContentType(fileBytes)
	if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/jpg" {
		return File{}, errors.New("incorrect file type")
	}

	img, _, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return File{}, err
	}

	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	if width < minWidth {
		return File{}, fmt.Errorf(`image width is %v less than %v`, width, minWidth)
	}

	if height < minHeight {
		return File{}, fmt.Errorf(`image height is %v less than %v`, height, minHeight)
	}

	if width > maxWidth {
		return File{}, fmt.Errorf(`image width is %v larger than %v`, width, maxWidth)
	}

	if height > maxHeight {
		return File{}, fmt.Errorf(`image height is %v larger than %v`, height, maxHeight)
	}

	if squareDimRequired && height != width {
		return File{}, fmt.Errorf(`image dimensions are not square: %v x %v`, strconv.Itoa(width), strconv.Itoa(height))
	}

	fileExt := strings.Split(filetype, "/")[1]
	return File{
		Name:  shortuuid.New() + "." + fileExt,
		Bytes: fileBytes,
	}, nil
}
