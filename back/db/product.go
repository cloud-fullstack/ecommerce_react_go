package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Product struct {
	ID                  string `json:"id"`
	Owner               string `json:"owner" binding:"required"`
	Store               string `json:"store" binding:"required"`
	Price               int    `json:"price" binding:"required,gte=100"`
	DiscountActive      bool   `json:"discount_active"`
	DiscountedPrice     int    `json:"discounted_price" binding:"required,ltefield=Price"`
	CreationDate        int    `json:"creation_date"`
	UpdateDate          int    `json:"update_date"`
	Listed              bool   `json:"listed"`
	SketchfabLink       string `json:"sketchfab_link"`
	YoutubeLink         string `json:"youtube_link"`
	TurntableLink       string `json:"turntable_link"`
	TurntableSlides     int    `json:"turntable_slides"`
	Name                string `json:"name" binding:"required"`
	Description         string `json:"description" binding:"required"`
	StoreName           string `json:"store_name"`
	BlogotexApplication string `json:"blogotex_application"`
	FAQEnabled          bool   `json:"faq_enabled"`
	QAEnabled           bool   `json:"qa_enabled"`
	Category            string `json:"category"`
}

type ProductInventoryLine struct {
	ProductID       string `json:"product_id"`
	InventoryItemID string `json:"inventory_item_id"`
	DemoItem        bool   `json:"demo_item"`
	Copyable        bool   `json:"copyable"` /* only for the handler */
}

type OrderProducts struct {
	OrderID string `json:"order_id"`
}

type ProductPicture struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Link      string `json:"link"`
	Index     int    `json:"index"`
}

type ProductPictureLove struct {
	PictureID     string `json:"picture_id"`
	LovedByAvatar bool   `json:"loved_by_avatar"`
	Loves         int    `json:"loves"`
}

type InventoryItemDropboxOwnerName struct {
	InventoryItem   DropboxInventoryItem `json:"inventory_item"`
	Dropbox         Dropbox              `json:"dropbox"`
	OwnerID         string               `json:"owner_id"`
	OwnerLegacyName string               `json:"owner_legacy_name"`
}

type FrontpageProductPreview struct {
	ProductID       string `json:"product_id"`
	ProductName     string `json:"product_name"`
	StoreID         string `json:"store_id"`
	StoreName       string `json:"store_name"`
	PictureLink     string `json:"picture_link"`
	Price           int    `json:"price"`
	Discounted      bool   `json:"discounted"`
	DiscountedPrice int    `json:"discounted_price"`
	Category        string `json:"category"`
}
type DiscountedProductsFrontpage struct {
	ProductID       string `json:"product_id"`
	ProductName     string `json:"product_name"`
	StoreID         string `json:"store_id"`
	StoreName       string `json:"store_name"`
	PictureLink     string `json:"picture_link"`
	Price           int    `json:"price"`
	Discounted      bool   `json:"discounted"`
	DiscountedPrice int    `json:"discounted_price"`
}

// BoughtProduct returns whether buyer has bought productID.
func BoughtProduct(ctx context.Context, conn Conn, productID string, buyer string) (bool, error) {
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return false, err
	}
	if buyer == "" {
		buyer = "00000000-0000-0000-0000-000000000000"
	}
	rows, err := tx.Query(ctx, `
		SELECT co.fulfilled, copl.demo FROM customer_order co
		JOIN customer_order_product_line copl ON co.id = copl.customer_order
		WHERE copl.product_id = $1 AND co.buyer = $2`,
		productID,
		buyer,
	)
	for rows.Next() {
		// todo: add if they are owner
		var fulfilled bool
		var demo bool
		err = rows.Scan(&fulfilled, &demo)
		if err != nil {
			return false, err
		}
		if fulfilled && !demo {
			return true, nil
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return false, err
	}
	return false, nil
}

func DeleteProduct(ctx context.Context, conn Conn, productID string) error {
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return nil
	}

	// deleting product: ERROR: update or delete on table "product" violates foreign key constraint "customer_order_product_line_product_id_fkey" on table "customer_order_product_line" (SQLSTATE 23503)

	_, err = tx.Exec(ctx, `
	DELETE
	FROM customer_order_product_line copl
	WHERE copl.product_id IN (
		SELECT
		copl.product_id
		FROM product p
		join customer_order_product_line copl ON p.id = copl.product_id
		WHERE p.id = $1
	)
	`, productID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
	DELETE FROM product_picture
	WHERE product_picture.product_id IN (
		SELECT
			product_picture.product_id
		FROM product
		JOIN product_picture pp ON product.id = pp.product_id
		WHERE product.id = $1
	)
	`, productID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
	DELETE FROM product_inventory_line
	WHERE product_inventory_line.product_id IN (
		SELECT
			pil.product_id
		FROM product_inventory_line pil
		JOIN product p ON pil.product_id = p.id
		WHERE p.id = $1
	)
	`, productID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM product WHERE id = $1", productID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil
	}

	return nil
}

func ProductOwner(ctx context.Context, conn Conn, productID string) (string, error) {
	var owner string
	err := conn.QueryRow(ctx, `
		SELECT owner
		FROM product
		WHERE id = $1`, productID).Scan(&owner)
	if err != nil {
		return "", err
	}
	return owner, nil
}

// FrontpageProductPreviews gets all products and randomizes the order.
func FrontpageProductPreviews(ctx context.Context, conn Conn) ([]FrontpageProductPreview, error) {
	rows, err := conn.Query(ctx, `
		SELECT
			p.id,
			p.name,
			s.id,
			s.name,
			i.link,
			p.price,
			p.discount_active,
			p.discounted_price,
			p.category
		FROM
			product p
			JOIN store s ON p.store = s.id
			JOIN product_picture i ON p.id = i.product_id
		WHERE
			p.listed = true AND
			i.index = 0
		ORDER BY
			random()
	`)
	if err != nil {
		return nil, err
	}
	var products []FrontpageProductPreview
	for rows.Next() {
		var product FrontpageProductPreview
		err = rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.StoreID,
			&product.StoreName,
			&product.PictureLink,
			&product.Price,
			&product.Discounted,
			&product.DiscountedPrice,
			&product.Category,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
func DiscountedProductsFrontpages(ctx context.Context, conn Conn) ([]DiscountedProductsFrontpage, error) {
	rows, err := conn.Query(ctx, `
		SELECT
			p.id,
			p.name,
			s.id,
			s.name,
			i.link,
			p.price,
			p.discount_active,
			p.discounted_price
		FROM
			product p
			JOIN store s ON p.store = s.id
			JOIN product_picture i ON p.id = i.product_id
		WHERE
			p.listed = true AND
			i.index = 0 AND 
			p.discount_active = true
		ORDER BY
			random()
	`)
	if err != nil {
		return nil, err
	}
	var products []DiscountedProductsFrontpage
	for rows.Next() {
		var product DiscountedProductsFrontpage
		err = rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.StoreID,
			&product.StoreName,
			&product.PictureLink,
			&product.Price,
			&product.Discounted,
			&product.DiscountedPrice,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
func ProductPictures(ctx context.Context, conn Conn, productID string) ([]ProductPicture, error) {
	rows, err := conn.Query(ctx, "SELECT * FROM product_picture WHERE product_id = $1", productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	pics := []ProductPicture{}
	for rows.Next() {
		var pic ProductPicture
		err = rows.Scan(&pic.ID, &pic.ProductID, &pic.Link, &pic.Index)
		if err != nil {
			return nil, err
		}
		pics = append(pics, pic)
	}
	return pics, nil
}

func ProductPicturesLoves(ctx context.Context, conn Conn, pictures []ProductPicture, avatarViewing string) ([]ProductPictureLove, error) {
	if avatarViewing == "" {
		avatarViewing = "00000000-0000-0000-0000-000000000000"
	}
	var picIDs []string
	for _, pic := range pictures {
		picIDs = append(picIDs, pic.ID)
	}
	ids := &pgtype.UUIDArray{}
	err := ids.Set(picIDs)
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, `
		SELECT product_picture.id,
		a.count,
		$1 IN (select ppl.owner from product_picture_love ppl where picture_id = product_picture.id) as liked_by_avatar
		FROM product_picture
		JOIN (
			SELECT product_picture.id as picture_id, COUNT(ppl.picture_id) as count
			FROM product_picture
					 LEFT JOIN product_picture_love ppl ON product_picture.id = ppl.picture_id
			GROUP BY product_picture.id
		) a ON a.picture_id = product_picture.id
		WHERE product_picture.id = ANY ($2)
	`, avatarViewing, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	loves := []ProductPictureLove{}
	for rows.Next() {
		var love ProductPictureLove
		err = rows.Scan(&love.PictureID, &love.Loves, &love.LovedByAvatar)
		if err != nil {
			return nil, err
		}
		loves = append(loves, love)
	}

	return loves, nil
}

func InsertProductPictureLove(ctx context.Context, conn Conn, pictureID, avatarID string) error {
	loved, err := ProductPictureLovedByAvatar(ctx, conn, pictureID, avatarID)
	if err != nil {
		return err
	}
	if loved {
		return errors.New("already loved by avatar")
	}
	_, err = conn.Exec(ctx, `INSERT INTO product_picture_love 
		VALUES(uuid_generate_v4(), $1, $2, to_timestamp($3))`,
		avatarID, pictureID, time.Now().Unix(),
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteProductPictureLove(ctx context.Context, conn Conn, pictureID, avatarID string) error {
	loved, err := ProductPictureLovedByAvatar(ctx, conn, pictureID, avatarID)
	if err != nil {
		return err
	}
	if !loved {
		return errors.New("picture is not loved by avatar")
	}

	_, err = conn.Exec(ctx, `DELETE FROM product_picture_love 
		WHERE $1 = picture_id AND owner = $2`, pictureID, avatarID,
	)
	if err != nil {
		return err
	}
	return nil
}

func ProductPictureLovedByAvatar(ctx context.Context, conn Conn, pictureID string, avatarID string) (bool, error) {
	var loved bool
	err := conn.QueryRow(ctx, `SELECT owner = $2 from 
		product_picture_love 
		WHERE picture_id = $1 AND owner = $2`, pictureID, avatarID).Scan(&loved)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
	}
	return loved, nil
}

func changeProductPictures(ctx context.Context, conn Conn, productID string, pics []ProductPicture) error {
	_, err := conn.Exec(ctx, "DELETE FROM product_picture WHERE product_id = $1", productID)
	for _, pic := range pics {
		_, err = conn.Exec(ctx, "INSERT INTO product_picture  VALUES (uuid_generate_v4(),$1, $2, $3)", productID, pic.Link, pic.Index)
		if err != nil {
			return err
		}
	}
	return nil
}

func NonCopyableInventoryItemsInAnotherProduct(ctx context.Context, conn *pgxpool.Pool, inventoryItemIDs []string, avatarKey, productID string) ([]DropboxInventoryItem, error) {
	ids := &pgtype.UUIDArray{}
	err := ids.Set(inventoryItemIDs)
	if productID == "" {
		productID = "00000000-0000-0000-0000-000000000000"
	}
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(ctx, `
	SELECT
    inventory_item.*
	FROM inventory_item
	join product_inventory_line pil ON inventory_item.id = pil.inventory_item_id
	join product p ON pil.product_id = p.id
	WHERE 
        inventory_item.id = ANY ($1) 
        AND NOT inventory_item.copyable
        AND p.owner = $2
		AND p.id <> $3`, ids, avatarKey, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var inventoryItems []DropboxInventoryItem
	for rows.Next() {
		var ii DropboxInventoryItem
		var t time.Time
		err = rows.Scan(
			&ii.ID,
			&ii.Name,
			&ii.ObjectID,
			&ii.DropboxID,
			&ii.Perms,
			&ii.Copyable,
			&t,
		)
		ii.AcquireTime = int(t.Unix())
		if err != nil {
			return nil, err
		}
		inventoryItems = append(inventoryItems, ii)
	}
	return inventoryItems, nil
}

func ProductsWithNoCopyItems(ctx context.Context, conn Conn, inventoryIDs []string) ([]Product, error) {
	ids := &pgtype.UUIDArray{}
	err := ids.Set(inventoryIDs)
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(ctx, `
	SELECT
		product.*
	FROM product
	JOIN product_inventory_line pil ON product.id = pil.product_id
	JOIN inventory_item ii ON pil.inventory_item_id = ii.id
	WHERE
		ii.id = ANY ($1)
		AND NOT ii.copyable`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var product Product
		creationDate := time.Time{}
		updateDate := time.Time{}
		err = rows.Scan(&product.ID, &product.Owner, &product.Store, &product.Price, &product.DiscountActive, &product.DiscountedPrice, &creationDate, &updateDate, &product.Listed, &product.SketchfabLink, &product.YoutubeLink, &product.TurntableLink, &product.TurntableSlides, &product.Name, &product.Description)
		product.CreationDate = int(creationDate.Unix())
		product.UpdateDate = int(updateDate.Unix())
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

// GetProductInventoryItemsAndDropboxURL gets a product's associated inventory items with their dropbox information.
func GetProductInventoryItemsAndDropboxURL(ctx context.Context, conn *pgxpool.Pool, productID string) ([]InventoryItemDropboxOwnerName, error) {
	var invItems []InventoryItemDropboxOwnerName
	rows, err := conn.Query(ctx, `
		SELECT inventory_item.*,
		dropbox.*, avatar.uuid, avatar.legacyname 
		FROM inventory_item
			JOIN product_inventory_line pil ON inventory_item.id = pil.inventory_item_id
			JOIN dropbox ON inventory_item.dropbox_id = dropbox.object_id
			JOIN product p ON pil.product_id = p.id
			JOIN avatar ON p.owner = avatar.uuid
		WHERE p.id = $1`, productID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		invItem := InventoryItemDropboxOwnerName{}
		var acquireTime time.Time
		err = rows.Scan(
			&invItem.InventoryItem.ID,
			&invItem.InventoryItem.Name,
			&invItem.InventoryItem.ObjectID,
			&invItem.InventoryItem.DropboxID,
			&invItem.InventoryItem.Perms,
			&invItem.InventoryItem.Copyable,
			&acquireTime,
			&invItem.Dropbox.ObjectID,
			nil,
			&invItem.Dropbox.Owner,
			&invItem.Dropbox.Region,
			&invItem.Dropbox.URL,
			&invItem.Dropbox.PosX,
			&invItem.Dropbox.PosY,
			&invItem.Dropbox.PosZ,
			&invItem.OwnerID,
			&invItem.OwnerLegacyName,
		)
		invItem.InventoryItem.AcquireTime = int(acquireTime.Unix())
		invItems = append(invItems, invItem)
		if err != nil {
			return nil, err
		}
	}
	return invItems, nil
}

func GetProductInventoryItems(ctx context.Context, conn *pgxpool.Pool, productID string) ([]DropboxInventoryItemWithDemo, error) {
	rows, err := conn.Query(ctx, `SELECT inventory_item.*, demo_item FROM inventory_item
    JOIN product_inventory_line pil ON inventory_item.id = pil.inventory_item_id
    JOIN product p ON pil.product_id = p.id
	WHERE p.id = $1`,
		productID)
	if err != nil {
		return []DropboxInventoryItemWithDemo{}, err
	}
	defer rows.Close()
	inventoryItems := make([]DropboxInventoryItemWithDemo, 0, 0)
	for rows.Next() {
		var inventoryItem DropboxInventoryItemWithDemo
		var t time.Time
		err = rows.Scan(&inventoryItem.ID, &inventoryItem.Name, &inventoryItem.ObjectID, &inventoryItem.DropboxID, &inventoryItem.Perms, &inventoryItem.Copyable, &t, &inventoryItem.DemoItem)
		if err != nil {
			return []DropboxInventoryItemWithDemo{}, err
		}
		inventoryItem.AcquireTime = int(t.Unix())
		inventoryItems = append(inventoryItems, inventoryItem)
	}
	return inventoryItems, nil
}

func GetAvatarProduct(ctx context.Context, conn *pgxpool.Pool, productID string) (Product, error) {
	product := Product{}
	creationDate := time.Time{}
	updateDate := time.Time{}
	err := conn.QueryRow(ctx, `
		SELECT product.*, store.name FROM product JOIN store ON product.store = store.id  WHERE product.id = $1 
	`, productID).Scan(&product.ID,
		&product.Owner,
		&product.Store,
		&product.Price,
		&product.DiscountActive,
		&product.DiscountedPrice,
		&creationDate,
		&updateDate,
		&product.Listed,
		&product.SketchfabLink,
		&product.YoutubeLink,
		&product.TurntableLink,
		&product.TurntableSlides,
		&product.Name,
		&product.Description,
		&product.BlogotexApplication,
		&product.FAQEnabled,
		&product.QAEnabled,
		&product.Category,
		&product.StoreName,
	)
	if err != nil {
		return Product{}, err
	}
	product.CreationDate = int(creationDate.Unix())
	product.UpdateDate = int(updateDate.Unix())
	if err != nil {
		return Product{}, err
	}
	return product, nil
}

func GetAvatarProducts(ctx context.Context, conn *pgxpool.Pool, avatarKey string) ([]Product, error) {
	rows, err := conn.Query(ctx, `SELECT * FROM product WHERE owner = $1`, avatarKey)
	if err != nil {
		return []Product{}, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		creationDate := time.Time{}
		updateDate := time.Time{}
		err = rows.Scan(&product.ID, &product.Owner, &product.Store, &product.Price, &product.DiscountActive, &product.DiscountedPrice, &creationDate, &updateDate, &product.Listed, &product.SketchfabLink, &product.YoutubeLink, &product.TurntableLink, &product.TurntableSlides, &product.Name, &product.Description, &product.BlogotexApplication, &product.FAQEnabled, &product.QAEnabled, &product.Category)
		if err != nil {
			return []Product{}, err
		}
		product.CreationDate = int(creationDate.Unix())
		product.UpdateDate = int(updateDate.Unix())
		products = append(products, product)
	}
	return products, nil
}

func UpsertProduct(ctx context.Context, conn *pgxpool.Pool, product Product, inventoryLines []ProductInventoryLine, pictures []ProductPicture, faqs []FAQ) (string, error) {
	idReturned := ""
	var err error = nil

	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	// if err != nil {
	// 	return "", err
	// }

	if product.ID == "" {
		err = tx.QueryRow(ctx,
			`INSERT INTO product 
             	VALUES(uuid_generate_v4(), $1, $2, $3, $4, $5, to_timestamp($6), to_timestamp($7), $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) RETURNING id`,
			product.Owner,
			product.Store,
			product.Price,
			product.DiscountActive,
			product.DiscountedPrice,
			product.CreationDate,
			product.UpdateDate,
			product.Listed,
			product.SketchfabLink,
			product.YoutubeLink,
			product.TurntableLink,
			product.TurntableSlides,
			product.Name,
			product.Description,
			product.BlogotexApplication,
			product.FAQEnabled,
			product.QAEnabled,
			product.Category,
		).Scan(&idReturned)
	} else {
		_, err = tx.Exec(ctx, `
					UPDATE product SET store = $2, price = $3, discount_active = $4, discounted_price = $5, creation_date = to_timestamp($6), 
						update_date = to_timestamp($7), listed = $8, sketchfab_link = $9, youtube_link = $10, turntable_link = $11, turntable_slides = $12,
						name = $13, description = $14, blogotex_application = $15, faq_enabled = $16, qa_enabled = $17, category = $18
					WHERE product.id = $1`,
			product.ID,
			product.Store,
			product.Price,
			product.DiscountActive,
			product.DiscountedPrice,
			product.CreationDate,
			product.UpdateDate,
			product.Listed,
			product.SketchfabLink,
			product.YoutubeLink,
			product.TurntableLink,
			product.TurntableSlides,
			product.Name,
			product.Description,
			product.BlogotexApplication,
			product.FAQEnabled,
			product.QAEnabled,
			product.Category,
		)
	}
	// if err != nil {
	// 	return "", err
	// }
	if idReturned == "" {
		idReturned = product.ID
	}
	err = upsertProductInventoryItems(ctx, tx, inventoryLines, idReturned)
	if err != nil {
		return "", err
	}

	err = changeProductPictures(ctx, tx, idReturned, pictures)
	if err != nil {
		return "", err
	}

	err = DeleteAllFAQs(ctx, tx, idReturned)
	if err != nil {
		return "", err
	}

	for _, faq := range faqs {
		err = InsertFAQ(ctx, tx, idReturned, faq.QuestionText, faq.AnswerText)
		if err != nil {
			return "", err
		}
	}

	err = tx.Commit(ctx)

	return idReturned, nil
}

func ProductsAssociatedWithDropbox(ctx context.Context, conn Conn, dropboxID string) ([]Product, error) {
	rows, err := conn.Query(ctx, `
		SELECT p.*
		FROM dropbox
			JOIN inventory_item ii ON dropbox.object_id = ii.dropbox_id
			JOIN product_inventory_line pil ON ii.id = pil.inventory_item_id
			JOIN product p ON pil.product_id = p.id
		WHERE dropbox.object_id = $1
		GROUP BY p.id`, dropboxID)
	if err != nil {
		return nil, err
	}

	products := make([]Product, 0, 0)
	//scan rows into products
	for rows.Next() {
		var product Product
		err = rows.Scan(
			&product.ID,
			&product.Owner,
			&product.Store,
			&product.Price,
			&product.DiscountActive,
			&product.DiscountedPrice,
			nil,
			nil,
			&product.Listed,
			&product.SketchfabLink,
			&product.YoutubeLink,
			&product.TurntableLink,
			&product.TurntableSlides,
			&product.Name,
			&product.Description,
			&product.BlogotexApplication,
			&product.FAQEnabled,
			&product.QAEnabled,
			&product.Category,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func ProductsAssociatedWithInventoryItem(ctx context.Context, conn Conn, inventoryItemID string) ([]Product, error) {
	rows, err := conn.Query(ctx, `
		SELECT product.*
		FROM product
			join product_inventory_line pil ON product.id = pil.product_id
			join inventory_item ii ON pil.inventory_item_id = ii.id
		WHERE ii.id = $1`, inventoryItemID)
	if err != nil {
		return nil, err
	}

	products := make([]Product, 0, 0)
	//scan rows into products
	for rows.Next() {
		var product Product
		err = rows.Scan(
			&product.ID,
			&product.Owner,
			&product.Store,
			&product.Price,
			&product.DiscountActive,
			&product.DiscountedPrice,
			nil,
			nil,
			&product.Listed,
			&product.SketchfabLink,
			&product.YoutubeLink,
			&product.TurntableLink,
			&product.TurntableSlides,
			&product.Name,
			&product.Description,
			&product.BlogotexApplication,
			&product.FAQEnabled,
			&product.QAEnabled,
			&product.Category,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func upsertProductInventoryItems(ctx context.Context, tx pgx.Tx, inventoryLines []ProductInventoryLine, productID string) error {
	for _, itemID := range inventoryLines {
		_, err := tx.Exec(ctx, `
			INSERT INTO product_inventory_line VALUES($1, $2, $3) ON CONFLICT ON CONSTRAINT unique_line DO NOTHING
		`, productID, itemID.InventoryItemID, itemID.DemoItem)
		if err != nil {
			return err
		}
	}

	err := deleteOrphanInventoryLines(ctx, tx, inventoryLines, productID)
	if err != nil {
		return err
	}

	return nil
}

func deleteOrphanInventoryLines(ctx context.Context, tx pgx.Tx, inventoryLines []ProductInventoryLine, productID string) error {
	rows, err := tx.Query(ctx, `SELECT * from product_inventory_line WHERE product_id = $1`,
		productID)
	if err != nil {
		return err
	}
	defer rows.Close()
	var productInventoryLines []ProductInventoryLine
	for rows.Next() {
		var productInventoryLine ProductInventoryLine
		err := rows.Scan(&productInventoryLine.ProductID, &productInventoryLine.InventoryItemID, &productInventoryLine.DemoItem)
		if err != nil {
			return err
		}
		productInventoryLines = append(productInventoryLines, productInventoryLine)
	}

	type Pair struct {
		String  string
		Boolean bool
	}
	inventoryItemIDsMap := make(map[Pair]struct{})
	for _, id := range inventoryLines {
		inventoryItemIDsMap[Pair{
			String:  id.InventoryItemID,
			Boolean: id.DemoItem,
		}] = struct{}{}
	}

	for _, productInventoryLine := range productInventoryLines {
		if _, exists := inventoryItemIDsMap[Pair{
			String:  productInventoryLine.InventoryItemID,
			Boolean: productInventoryLine.DemoItem,
		}]; !exists {
			_, err := tx.Exec(ctx, `DELETE FROM product_inventory_line WHERE product_id = $1 AND inventory_item_id = $2 AND demo_item = $3`, productID, productInventoryLine.InventoryItemID, productInventoryLine.DemoItem)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func UnlistProducts(ctx context.Context, conn Conn, productIDs []string) error {
	// todo: make hud delete itself after taking an order
	// todo: check that a product has items after being deleted normally
	ids := &pgtype.UUIDArray{}
	err := ids.Set(productIDs)
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, `UPDATE product SET listed = false WHERE id = ANY ($1)`, ids)
	if err != nil {
		return err
	}
	return nil
}
