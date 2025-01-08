package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"net/http"
	"shopa/db"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Dropbox struct {
	ObjectID string  `json:"object_id"`
	Owner    string  `json:"owner"`
	Region   string  `json:"region"`
	URL      string  `json:"url"`
	Online   bool    `json:"online"`
	PosX     float32 `json:"pos_x"`
	PosY     float32 `json:"pos_y"`
	PosZ     float32 `json:"pos_z"`
}

// todo: verify all error messages are proper and have returns in the proper areas and proper http codes, and logs for successes

func DeliverDropbox(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	requesterDetails := struct {
		RequesterAvatarKey  string
		RequesterAvatarName string
		Region              string
		ObjectKey           string
		ObjectName          string
	}{
		RequesterAvatarKey:  c.GetHeader("x-secondlife-owner-key"),
		RequesterAvatarName: c.GetHeader("x-secondlife-owner-name"),
		Region:              c.GetHeader("x-secondlife-region"),
		ObjectKey:           c.GetHeader("x-secondlife-object-key"),
		ObjectName:          c.GetHeader("x-secondlife-object-name"),
	}

	repos, err := db.GetDropboxRepos(c.Request.Context(), dbConn)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP":               c.ClientIP(),
			"requesterDetails": requesterDetails,
		}, 500, "error getting dropbox repos: "+err.Error())
		return
	}

	repoURLS := make([]string, 0, len(repos))
	for _, dropbox := range repos {
		repoURLS = append(repoURLS, dropbox.URL)
	}
	dropboxStatuses := pingDropboxURLs(repoURLS)
	onlineURL := ""
	for url, online := range dropboxStatuses {
		if online {
			onlineURL = url
			break
		}
	}
	if onlineURL == "" {
		logRespondError(c, log.Fields{
			"IP":               c.ClientIP(),
			"requesterDetails": requesterDetails,
		}, 400, "no dropbox repos online")
		return
	}

	respMessageJSON := struct {
		Message string `json:"message"`
		Sent    bool   `json:"sent"`
	}{}
	reqBody := fmt.Sprintf(`{"message":"dropbox_delivery", "requester":"%s"}`, requesterDetails.RequesterAvatarKey)
	resp, err := http.Post(onlineURL, "application/json", strings.NewReader(reqBody))
	if err != nil {
		logRespondError(c, log.Fields{
			"IP":               c.ClientIP(),
			"requesterDetails": requesterDetails,
		}, 500, "sending POST delivery request to dropbox: "+err.Error())
		return
	}
	r, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(r, &respMessageJSON)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP":               c.ClientIP(),
			"requesterDetails": requesterDetails,
		}, 500, "reading response from dropbox repo: "+err.Error())
		return
	}

	if respMessageJSON.Sent != true {
		logRespondError(c, log.Fields{
			"IP":               c.ClientIP(),
			"requesterDetails": requesterDetails,
			"dropboxResponse":  respMessageJSON,
		}, 500, "dropbox repo unable to send dropbox")
		return
	}

	log.WithFields(log.Fields{
		"IP":               c.ClientIP(),
		"requesterDetails": requesterDetails,
		"dropboxResponse":  respMessageJSON,
	}).Info("delivered dropbox")
	c.JSON(200, gin.H{
		"message": "delivered dropbox",
	})
}

func GetAvatarDropboxContents(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)

	userInfo := struct {
		DropboxID string `json:"dropbox_ID"`
	}{}
	err := c.ShouldBindJSON(&userInfo)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding userInfo:"+err.Error())
		return
	}

	dropboxes, err := db.GetAvatarDropboxContents(c.Request.Context(), dbConn, userInfo.DropboxID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "getting dropbox contents from db:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":              c.ClientIP(),
		"dropboxContents": dropboxes,
	}).Info("executed GetAvatarDropboxesContents")
	c.JSON(200, dropboxes)
}

func GetAvatarDropboxes(c *gin.Context) {
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

	dropboxes, err := db.GetAvatarDropboxes(c.Request.Context(), dbConn, userInfo.AvatarKey)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "getting avatar's dropboxes from db:"+err.Error())
		return
	}

	dropboxURLs := make([]string, 0, len(dropboxes))
	for _, dropbox := range dropboxes {
		dropboxURLs = append(dropboxURLs, dropbox.URL)
	}

	// todo: remove if the dropbox status is not online
	// todo: add an unlisting if not online or missing
	// todo: hide URL from heartbeating the HUD in pingHUD
	dropboxStatuses := pingDropboxURLs(dropboxURLs)

	dropboxesWithStatus := make([]Dropbox, 0, len(dropboxes))
	for _, dropbox := range dropboxes {
		if dropboxStatuses[dropbox.URL] == false {
			continue
		}
		dropboxesWithStatus = append(dropboxesWithStatus, Dropbox{
			ObjectID: dropbox.ObjectID,
			Owner:    dropbox.Owner,
			Region:   dropbox.Region,
			URL:      dropbox.URL,
			Online:   dropboxStatuses[dropbox.URL],
			PosX:     dropbox.PosX,
			PosY:     dropbox.PosY,
			PosZ:     dropbox.PosZ,
		})

	}

	log.WithFields(log.Fields{
		"IP":                   c.ClientIP(),
		"DropboxInventoryItem": dropboxes,
	}).Info("executed GetAvatarDropboxes")
	c.JSON(200, dropboxesWithStatus)
}

func pingDropboxURLs(urls []string) map[string]bool {
	wg := sync.WaitGroup{}
	client := http.Client{
		Timeout: time.Second * 30,
	}

	urlPingStatuses := make(map[string]bool)
	mu := sync.Mutex{}

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			success, _ := pingDropboxURL(url, &client)
			mu.Lock()
			urlPingStatuses[url] = success
			mu.Unlock()
		}(url)
	}
	wg.Wait()
	return urlPingStatuses
}

func pingDropboxURL(url string, client *http.Client) (bool, error) {
	r := strings.NewReader(`{"message":"ping"}`)
	resp, err := client.Post(url, "application/json", r)
	if err != nil {
		return false, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	expectedResp := struct {
		Status string `json:"status"`
	}{}

	err = json.Unmarshal(b, &expectedResp)
	if err != nil {
		return false, err
	}

	if expectedResp.Status != "ok" {
		return false, errors.New("did not get expected \"ok\" from dropbox ")
	}

	return true, nil
}

func UpdateDropboxContents(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	dropboxInventoryItems := make([]db.DropboxInventoryItem, 0, 0)
	err := c.ShouldBindJSON(&dropboxInventoryItems)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding dropboxInventoryItems:"+err.Error())
		return
	}

	orphanInventoryItems, productsWithOrphanedItems, err := db.UpdateDropboxContents(c.Request.Context(), dbConn, dropboxInventoryItems)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP":                    c.ClientIP(),
			"dropboxInventoryItems": dropboxInventoryItems,
		}, 500, "updating dropbox contents to db: "+err.Error())
		return
	}

	if len(productsWithOrphanedItems) > 0 {
		var productIDs []string
		for _, product := range productsWithOrphanedItems {
			productIDs = append(productIDs, product.ID)
		}
		err = db.UnlistProducts(c.Request.Context(),dbConn, productIDs)
		if err != nil {
			logRespondError(c, log.Fields{
				"IP":                    c.ClientIP(),
				"dropboxInventoryItems": dropboxInventoryItems,
			}, 500, "unlisting products with orphan items: " + err.Error())
			return
		}
		var productOwner string
		var productNames []string
		for _, p := range productsWithOrphanedItems {
			productOwner = p.Owner
			productNames = append(productNames, p.Name)
		}
		var itemNames []string
		for _, item := range orphanInventoryItems {
			itemNames = append(itemNames, item.Name)
		}
		err = db.CreateNotification(
			c.Request.Context(),
			dbConn,
			productOwner,
			fmt.Sprintf(`The following products have become unlisted due to their items (%s) being deleted or sold: %v`, strings.Join(itemNames,","), strings.Join(productNames,",")),
		)
		if err != nil {
			logRespondError(c, log.Fields{
				"IP":                    c.ClientIP(),
				"dropboxInventoryItems": dropboxInventoryItems,
			}, 500, "creating notification for unlisted items: " + err.Error())
			return
		}
	}

	log.WithFields(log.Fields{
		"IP":                    c.ClientIP(),
		"dropboxInventoryItems": dropboxInventoryItems,
	}).Info("executed UpdateDropboxContents")
	c.JSON(200, gin.H{
		"message": "updated dropbox inventory items",
	})
}

func InsertDropbox(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)

	SLDropbox := struct {
		ObjectID string `json:"object_id"`
		Owner    string `json:"owner"`
		Region   string `json:"region"`
		URL      string `json:"url"`
		PosX     string `json:"pos_x"`
		PosY     string `json:"pos_y"`
		PosZ     string `json:"pos_z"`
	}{}
	err := c.ShouldBindJSON(&SLDropbox)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding SLDropbox:"+err.Error())
		return
	}

	stringToFloatConversionErrorReporting := func() {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "converting string to float:"+err.Error())
	}

	x, err := strconv.ParseFloat(SLDropbox.PosX, 32)
	if err != nil {
		stringToFloatConversionErrorReporting()
		return
	}
	y, err := strconv.ParseFloat(SLDropbox.PosY, 32)
	if err != nil {
		stringToFloatConversionErrorReporting()
		return
	}
	z, err := strconv.ParseFloat(SLDropbox.PosZ, 32)
	if err != nil {
		stringToFloatConversionErrorReporting()
		return
	}

	dropboxInfo := db.Dropbox{
		ObjectID: SLDropbox.ObjectID,
		Owner:    SLDropbox.Owner,
		Region:   SLDropbox.Region,
		URL:      SLDropbox.URL,
		PosX:     float32(x),
		PosY:     float32(y),
		PosZ:     float32(z),
	}

	err = db.InsertDropbox(c.Request.Context(), dbConn, dropboxInfo)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			log.WithFields(log.Fields{
				"IP":          c.ClientIP(),
				"dropBoxInfo": dropboxInfo,
			}).Info("dropbox already inserted")
			c.JSON(200, gin.H{
				"message": "dropbox already inserted",
			})
			return
		}
		if strings.Contains(err.Error(), `foreign key constraint "dropbox_owner_fkey"`) {
			log.WithFields(log.Fields{
				"IP":          c.ClientIP(),
				"dropBoxInfo": dropboxInfo,
			}).Info("dropbox owner does not exist in database")
			c.JSON(200, gin.H{
				"error":   true,
				"message": "Must wear HUD to register before rezzing dropbox.",
			})
			return
		}

		logRespondError(c, log.Fields{
			"IP":          c.ClientIP(),
			"dropBoxInfo": dropboxInfo,
		}, 500, "error inserting dropboxInfo: "+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":          c.ClientIP(),
		"dropBoxInfo": dropboxInfo,
	}).Info("inserted/updated dropbox")
	c.JSON(200, gin.H{
		"message": "inserted/updated dropbox",
	})
}

func InsertDropboxRepo(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)

	SLDropbox := struct {
		ObjectID string `json:"object_id"`
		Owner    string `json:"owner"`
		Region   string `json:"region"`
		URL      string `json:"url"`
		PosX     string `json:"pos_x"`
		PosY     string `json:"pos_y"`
		PosZ     string `json:"pos_z"`
	}{}
	err := c.ShouldBindJSON(&SLDropbox)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding json:"+err.Error())
		return
	}

	stringToFloatConversionErrorReporting := func() {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "converting string to float:"+err.Error())
	}

	x, err := strconv.ParseFloat(SLDropbox.PosX, 32)
	if err != nil {
		stringToFloatConversionErrorReporting()
		return
	}
	y, err := strconv.ParseFloat(SLDropbox.PosY, 32)
	if err != nil {
		stringToFloatConversionErrorReporting()
		return
	}
	z, err := strconv.ParseFloat(SLDropbox.PosZ, 32)
	if err != nil {
		stringToFloatConversionErrorReporting()
		return
	}

	dropboxInfo := db.Dropbox{
		ObjectID: SLDropbox.ObjectID,
		Owner:    SLDropbox.Owner,
		Region:   SLDropbox.Region,
		URL:      SLDropbox.URL,
		PosX:     float32(x),
		PosY:     float32(y),
		PosZ:     float32(z),
	}

	err = db.InsertDropboxRepo(c.Request.Context(), dbConn, dropboxInfo)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			log.WithFields(log.Fields{
				"IP":          c.ClientIP(),
				"dropBoxInfo": dropboxInfo,
			}).Info("dropbox already inserted")
			c.JSON(200, gin.H{
				"message": "dropbox already inserted",
			})
			return
		}
		if strings.Contains(err.Error(), `foreign key constraint "dropbox_owner_fkey"`) {
			log.WithFields(log.Fields{
				"IP":          c.ClientIP(),
				"dropBoxInfo": dropboxInfo,
			}).Error("dropbox owner does not exist in database")
			c.JSON(200, gin.H{
				"error":   true,
				"message": "Must wear HUD to register before rezzing dropbox.",
			})
			return
		}

		logRespondError(c, log.Fields{
			"IP":          c.ClientIP(),
			"dropBoxInfo": dropboxInfo,
		}, 500, "inserting dropbox repo:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":          c.ClientIP(),
		"dropBoxInfo": dropboxInfo,
	}).Info("inserted/updated dropbox repo")
	c.JSON(200, gin.H{
		"message": "inserted/updated dropbox repo",
	})
}
