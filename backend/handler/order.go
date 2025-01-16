package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/avast/retry-go"
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

type OrderInfo struct {
	OrderLines  []db.OrderLine `json:"order_lines"`
	AvatarBuyer string         `json:"avatar_buyer"`
}

type MissingErrors struct {
	Errors []error
}

func (m MissingErrors) Error() string {
	var errStrings []string
	for _, err := range m.Errors {
		errStrings = append(errStrings, err.Error())
	}
	return "multiple errors for missing entities: " + strings.Join(errStrings, ",")
}

func Dedupe(m MissingErrors) error {
	errMap := make(map[string]error)
	for _, err := range m.Errors {
		errMap[err.Error()] = err
	}
	// turn errMap into an array of errors
	var dedupedErrors []error
	for _, err := range errMap {
		dedupedErrors = append(dedupedErrors, err)
	}
	if len(dedupedErrors) == 1 {
		return dedupedErrors[0]
	}
	m.Errors = dedupedErrors
	return m
}

type MissingItemError struct {
	Item db.InventoryItemDropboxOwnerName
}

func (m MissingItemError) Error() string {
	return "missing item with the ID: " + m.Item.InventoryItem.ID
}

type MissingDropboxError struct {
	Dropbox db.Dropbox
}

func (m MissingDropboxError) Error() string {
	return "missing dropbox with the ID: " + m.Dropbox.ObjectID
}

type NoItemsError struct {
}

func (m NoItemsError) Error() string {
	return "there are no items in this order. Products may be under construction or unlisted"
}

// todo: disallow redelivery of gacha items

// ResendOrder resends an order to a HUD.
func ResendOrder(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	orderInfo := struct {
		OrderID     string `json:"order_id"`
		AvatarBuyer string `json:"avatar_buyer"`
	}{}
	err := c.ShouldBindJSON(&orderInfo)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding orderInfo:"+err.Error())
		return
	}

	successful, err := pingHUD(c.Request.Context(), dbConn, orderInfo.AvatarBuyer)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 400, "Unable to ping HUD: "+err.Error())
		return
	}

	if !successful {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 400, "Unable to ping HUD: unreachable")
		return
	}

	invoiceLines, err := db.OrderPayables(c.Request.Context(), dbConn, orderInfo.OrderID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 400, "creating order in DB: "+err.Error())
		return
	}

	err = sendOrderToHUD(c.Request.Context(), dbConn, invoiceLines, orderInfo.AvatarBuyer, orderInfo.OrderID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 400, "sending order to HUD: "+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":         c.ClientIP(),
		"order_info": orderInfo,
	}).Info("resent order to HUD")
	c.JSON(200, gin.H{
		"order_id": orderInfo.OrderID,
		"message":  "resent order to HUD",
	})
}

func RedeliverOrderProduct(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	redeliverInfo := struct {
		OrderID      string `json:"order_id"`
		ProductID    string `json:"product_id"`
		TargetAvatar string `json:"target_avatar"`
	}{}
	err := c.ShouldBindJSON(&redeliverInfo)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding redeliverInfo:"+err.Error())
		return
	}

	orderItems, err := db.ProductOrderItems(c.Request.Context(), dbConn, redeliverInfo.OrderID, redeliverInfo.ProductID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "getting orderItems:"+err.Error())
		return
	}

	for _, item := range orderItems {
		if !item.Copyable {
			logRespondError(c, log.Fields{
				"IP": c.ClientIP(),
			}, 500, "this product contains limited-inventory, no-copy items")
			return
		}
	}

	if len(orderItems) == 0 {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "no items to redeliver")
		return
	}

	orderItemsMap := make(map[string][]db.OrderItem)
	for _, orderItem := range orderItems {
		if _, exists := orderItemsMap[orderItem.DropboxURL]; !exists {
			orderItemsMap[orderItem.DropboxURL] = make([]db.OrderItem, 0, 0)
		}
		slice := orderItemsMap[orderItem.DropboxURL]
		slice = append(slice, orderItem)
		orderItemsMap[orderItem.DropboxURL] = slice
	}

	wg := sync.WaitGroup{}
	errorss := make([]error, 0, 0)

	for dropboxURL, orderItemsArr := range orderItemsMap {
		wg.Add(1)
		go func(dropboxURL string, orderItems []db.OrderItem) {
			err := commandDropboxSend(dropboxURL, orderItems, redeliverInfo.TargetAvatar, redeliverInfo.OrderID)
			if err != nil {
				errorss = append(errorss, err)
			}
			wg.Done()
		}(dropboxURL, orderItemsArr)
	}
	wg.Wait()

	if len(errorss) > 0 {
		var errStrings []string
		for _, err2 := range errorss {
			errStrings = append(errStrings, err2.Error())
		}
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "commanding dropboxes to send:"+strings.Join(errStrings, ","))
		return
	}

	log.WithFields(log.Fields{
		"IP":             c.ClientIP(),
		"redeliver_info": redeliverInfo.OrderID,
	}).Info("redelivered order product")
	c.JSON(200, gin.H{
		"message":        "redelivered order product",
		"redeliver_info": redeliverInfo.OrderID,
	})
}

func CompleteOrder(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	orderInfo := struct {
		OrderID      string `json:"order_id"`
		TargetAvatar string `json:"target_avatar"`
	}{}
	err := c.ShouldBindJSON(&orderInfo)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding orderInfo:"+err.Error())
		return
	}

	orderItems, err := db.OrderItems(c.Request.Context(), dbConn, orderInfo.OrderID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "getting orderItems:"+err.Error())
		return
	}

	orderItemsMap := make(map[string][]db.OrderItem)
	for _, orderItem := range orderItems {
		if _, exists := orderItemsMap[orderItem.DropboxURL]; !exists {
			orderItemsMap[orderItem.DropboxURL] = make([]db.OrderItem, 0, 0)
		}
		slice := orderItemsMap[orderItem.DropboxURL]
		slice = append(slice, orderItem)
		orderItemsMap[orderItem.DropboxURL] = slice
	}

	wg := sync.WaitGroup{}
	errorss := make([]error, 0, 0)

	for dropboxURL, orderItemsArr := range orderItemsMap {
		wg.Add(1)
		go func(dropboxURL string, orderItems []db.OrderItem) {
			err := commandDropboxSend(dropboxURL, orderItems, orderInfo.TargetAvatar, orderInfo.OrderID)
			if err != nil {
				errorss = append(errorss, err)
			}
			wg.Done()
		}(dropboxURL, orderItemsArr)
	}
	wg.Wait()

	if len(errorss) > 0 {
		var errStrings []string
		for _, err2 := range errorss {
			errStrings = append(errStrings, err2.Error())
		}
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "commanding dropboxes to send:"+strings.Join(errStrings, ","))
		return
	}

	err = db.MarkOrderFulfilled(c.Request.Context(), dbConn, orderInfo.OrderID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "completing order:"+err.Error())
		return
	}

	var noCopyInventoryDatabaseIDs []string
	for _, item := range orderItems {
		if !item.Copyable {
			noCopyInventoryDatabaseIDs = append(noCopyInventoryDatabaseIDs, item.DatabaseID)
		}
	}

	err = cleanupDBWithNoCopyItems(c.Request.Context(), dbConn, noCopyInventoryDatabaseIDs)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "completing order:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":       c.ClientIP(),
		"order_id": orderInfo.OrderID,
	}).Info("completed order")
	c.JSON(200, gin.H{
		"message":  "completed order",
		"order_id": orderInfo.OrderID,
	})
}

// cleanupDBWithNoCopyItems delists products, deletes inventory items pertaining to those in inventoryIDs that have no copy items.
// It additionally sends a notification to the avatarID.
func cleanupDBWithNoCopyItems(ctx context.Context, conn db.Conn, inventoryIDs []string) error {
	// Get products containing nocopy items
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return nil
	}
	products, err := db.ProductsWithNoCopyItems(ctx, tx, inventoryIDs)
	if err != nil {
		return err
	}

	var productIDs []string
	for _, product := range products {
		productIDs = append(productIDs, product.ID)
	}

	err = db.UnlistProducts(ctx, tx, productIDs)

	for _, iid := range inventoryIDs {
		err := db.DeleteInventoryItem(ctx, tx, iid)
		if err != nil {
			return err
		}
	}

	for _, product := range products {
		err = db.CreateNotification(ctx, tx, product.Owner, "Unlisted the following products for selling a no-copy item: "+product.Name)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func CancelOrder(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	orderInfo := struct {
		OrderID string `json:"order_id"`
	}{}
	err := c.ShouldBindJSON(&orderInfo)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding orderInfo:"+err.Error())
		return
	}

	err = db.CancelOrder(c.Request.Context(), dbConn, orderInfo.OrderID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "cancelling order:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":       c.ClientIP(),
		"order_id": orderInfo.OrderID,
	}).Info("cancelled order")
	c.JSON(200, gin.H{
		"message":  "cancelled order",
		"order_id": orderInfo.OrderID,
	})
}

func OrderHistory(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	info := struct {
		AvatarKey string `json:"avatar_key"`
	}{}
	err := c.ShouldBindJSON(&info)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding info:"+err.Error())
		return
	}
	orders, err := db.OrderHistory(c.Request.Context(), dbConn, info.AvatarKey)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "getting order history:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":   c.ClientIP(),
		"info": info,
	}).Info("completed order")
	c.JSON(200, orders)
}

// CreateOrder creates an order and sends it to the user's HUD.
// Also checks the dropboxes that have inventory items associated with the order and verifies they are there.
func CreateOrder(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	orderInfo := OrderInfo{}
	err := c.ShouldBindJSON(&orderInfo)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding orderInfo:"+err.Error())
		return
	}

	var productIDs []string
	for _, line := range orderInfo.OrderLines {
		productIDs = append(productIDs, line.ProductID)
	}

	for _, productID := range productIDs {
		product, err := db.GetAvatarProduct(c.Request.Context(), dbConn, productID)
		if err != nil {
			logRespondError(c, log.Fields{
				"IP": c.ClientIP(),
			}, 500, "getting product:"+err.Error())
			return
		}
		if !product.Listed {
			logRespondError(c, log.Fields{
				"IP": c.ClientIP(),
			}, 500, "the following product is unlisted and cannot be ordered: " + product.Name)
			return
		}
	}

	err = getMissingOrderItems(c.Request.Context(), dbConn, productIDs)
	if err != nil {
		var serros []error

		switch v := err.(type) {
		case MissingErrors:
			serros = v.Errors
		case MissingDropboxError:
			serros = append(serros, v)
		case MissingItemError:
			serros = append(serros, v)
		default:
			serros = append(serros, v)
		}

		var avatarOwnerID string
		var missingDropboxIDs []string
		var missingItemIDs []string

		//unusualError := false

		for _, err := range serros {
			switch serr := err.(type) {
			case MissingDropboxError:
				avatarOwnerID = serr.Dropbox.Owner
				missingDropboxIDs = append(missingDropboxIDs, serr.Dropbox.ObjectID)
			case MissingItemError:
				avatarOwnerID = serr.Item.OwnerID
				missingItemIDs = append(missingItemIDs, serr.Item.InventoryItem.ID)
			default:
				if strings.Contains(err.Error(), "context deadline") {
					logRespondError(c, log.Fields{
						"IP":         c.ClientIP(),
						"order_info": orderInfo,
					}, 500, "dropbox took too long to respond in creating order. Try again.")
					return
				}
				logRespondError(c, log.Fields{
					"IP":         c.ClientIP(),
					"order_info": orderInfo,
				}, 500, err.Error())
				return
			}

		}

		productsAffected, err := deleteEntities(c.Request.Context(), dbConn, missingDropboxIDs, missingItemIDs)
		if err != nil {
			logRespondError(c, log.Fields{
				"IP":         c.ClientIP(),
				"order_info": orderInfo,
			}, 500, "deleting missing entities:"+err.Error())
			return
		}

		var productIDsAffected []string
		for _, product := range productsAffected {
			productIDsAffected = append(productIDsAffected, product.ID)
		}
		err = db.UnlistProducts(c.Request.Context(), dbConn, productIDsAffected)
		if err != nil {
			logRespondError(c, log.Fields{
				"IP":         c.ClientIP(),
				"order_info": orderInfo,
			}, 500, "unlisting a product with missing items:"+err.Error())
			return
		}

		var productNamesAffected []string
		for _, product := range productsAffected {
			productNamesAffected = append(productNamesAffected, product.Name)
		}
		errCr := db.CreateNotification(c.Request.Context(), dbConn, avatarOwnerID,
			"The following products have become unlisted due to them missing an item in-world, or a dropbox containing their items missing in-world: "+strings.Join(productNamesAffected, ","))
		if errCr != nil {
			logRespondError(c, log.Fields{
				"IP":         c.ClientIP(),
				"order_info": orderInfo,
			}, 500, "creating a notification for unlisting products with missing items:"+err.Error())
			return
		}

		errMsg := fmt.Sprintf(`Your order contained products (%v) with items that no longer exist. Please try again after removing the products.`,
			strings.Join(productNamesAffected,","))
		logRespondError(c, log.Fields{
			"IP":     c.ClientIP(),
			"error": errMsg,
		}, 500, errMsg)
		return
	}

	successful, err := pingHUD(c.Request.Context(), dbConn, orderInfo.AvatarBuyer)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 400, "Unable to ping HUD: "+err.Error())
		return
	}
	if !successful {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 400, "Unable to ping HUD: unreachable")
		return
	}

	orderID, err := db.CreateOrder(c.Request.Context(), dbConn, orderInfo.OrderLines, orderInfo.AvatarBuyer)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 400, "creating order: "+err.Error())
		return
	}

	invoiceLines, err := db.OrderPayables(c.Request.Context(), dbConn, orderID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 400, "creating order in DB: "+err.Error())
		return
	}

	err = sendOrderToHUD(c.Request.Context(), dbConn, invoiceLines, orderInfo.AvatarBuyer, orderID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 400, "sending order to HUD: "+err.Error())
		return
	}

	//todo: unlist products that have no items listed
	// todo: add ability to delete stores and products

	log.WithFields(log.Fields{
		"IP":       c.ClientIP(),
		"order_id": orderID,
	}).Info("submitted order and sent it to avatar's HUD")
	c.JSON(200, gin.H{
		"message":  "submitted order and sent it to avatar's HUD",
		"order_id": orderID,
	})
}

// deleteEntities deletes dropboxIDs and missingItemIDs and returns the productIDs affected.
func deleteEntities(ctx context.Context, dbConn db.Conn, missingDropboxIDs []string, missingItemIDs []string) (productsAffected []db.Product, err error) {
	tx, err := dbConn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return nil, err
	}

	for _, dropboxID := range missingDropboxIDs {
		err = db.DeleteDropbox(ctx, tx, dropboxID)
		if err != nil {
			return nil, err
		}
	}

	for _, itemID := range missingItemIDs {
		err = db.DeleteInventoryItem(ctx, tx, itemID)
		if err != nil {
			return nil, err
		}
	}

	allProducts := make(map[db.Product]struct{})
	for _, dropboxID := range missingDropboxIDs {
		products, err := db.ProductsAssociatedWithDropbox(ctx, dbConn, dropboxID)
		if err != nil {
			return nil, err
		}
		for _, product := range products {
			allProducts[product] = struct{}{}
		}
	}

	for _, itemID := range missingItemIDs {
		products, err := db.ProductsAssociatedWithInventoryItem(ctx, dbConn, itemID)
		if err != nil {
			return nil, err
		}
		for _, product := range products {
			allProducts[product] = struct{}{}
		}
	}

	// convert allProducts to a slice
	productsSlice := make([]db.Product, 0, len(allProducts))
	for product := range allProducts {
		productsSlice = append(productsSlice, product)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return productsSlice, nil
}

func sendOrderToHUD(ctx context.Context, dbConn *pgxpool.Pool, invoiceLines []db.InvoiceLine, targetAvatar string, orderID string) error {
	var slDigestibleList []string
	for _, line := range invoiceLines {
		slDigestibleList = append(slDigestibleList, line.SellerAvatarUUID, line.SellerAvatarName, strconv.Itoa(line.ProductsSum))
	}
	row := dbConn.QueryRow(ctx, "SELECT hud_url FROM avatar WHERE uuid=$1", targetAvatar)
	hudURL := ""
	err := row.Scan(&hudURL)
	if err != nil {
		return err
	}

	b, err := json.Marshal(&struct {
		Message      string `json:"message"`
		InvoiceLines string `json:"invoice_lines"`
		OrderID      string `json:"order_id"`
	}{
		Message:      "order",
		InvoiceLines: strings.Join(slDigestibleList, ","),
		OrderID:      orderID,
	})
	if err != nil {
		return err
	}

	client := http.Client{Timeout: time.Second * 5}
	resp, err := client.Post(hudURL, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}
	if string(respBody) != "received" {
		return errors.New("did not get expected body after sending order to HUD: " + string(respBody))
	}
	return nil
}

func getMissingOrderItems(ctx context.Context, dbConn *pgxpool.Pool, productIDs []string) error {
	// need to get the owners of the items
	allErrors := make([]error, 0, 0)
	for _, productID := range productIDs {
		err := getMissingProductItems(ctx, dbConn, productID)
		if err != nil {
			allErrors = append(allErrors, err)
		}
	}

	if len(allErrors) == 1 {
		return allErrors[0]
	}

	if len(allErrors) > 1 {
		var parentError MissingErrors
		for _, err := range allErrors {
			if thisErr, ok := err.(*MissingErrors); ok {
				for _, nestedError := range thisErr.Errors {
					parentError.Errors = append(parentError.Errors, nestedError)
				}
				continue
			}
			parentError.Errors = append(parentError.Errors, err)
		}
		return Dedupe(parentError)
	}

	return nil
}

// getMissingProductItems gets the inventory items missing from all the dropboxes that contain inventory items associated with a productID.
func getMissingProductItems(ctx context.Context, dbConn *pgxpool.Pool, productID string) error {
	// todo: get dropbox name to be used in error message and collect missing items.
	items, err := db.GetProductInventoryItemsAndDropboxURL(ctx, dbConn, productID)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return NoItemsError{}
	}

	//[dropboxURL][]db.InventoryItemDropboxURLOwnerName
	itemsMap := make(map[db.Dropbox][]db.InventoryItemDropboxOwnerName)
	for _, item := range items {
		if _, exists := itemsMap[item.Dropbox]; !exists {
			itemsMap[item.Dropbox] = make([]db.InventoryItemDropboxOwnerName, 0, 0)
		}
		slice := itemsMap[item.Dropbox]
		slice = append(slice, item)
		itemsMap[item.Dropbox] = slice
	}

	// do some go routine stuff with getMissingDropboxItems
	wg := sync.WaitGroup{}
	allErrors := make([]error, 0, 0)
	mu := sync.Mutex{}

	for dropboxURL, itemsArr := range itemsMap {
		wg.Add(1)
		go func(dropbox db.Dropbox, items []db.InventoryItemDropboxOwnerName) {
			defer wg.Done()
			err := getMissingDropboxItems(dropbox, items)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				allErrors = append(allErrors, err)
			}
		}(dropboxURL, itemsArr)
	}
	wg.Wait()

	if len(allErrors) == 1 {
		return allErrors[0]
	}
	if len(allErrors) > 1 {
		var parentError MissingErrors
		for _, err := range allErrors {
			if thisErr, ok := err.(*MissingErrors); ok {
				for _, nestedError := range thisErr.Errors {
					parentError.Errors = append(parentError.Errors, nestedError)
				}
				continue
			}
			parentError.Errors = append(parentError.Errors, err)
		}
		return Dedupe(parentError)
	}

	return nil
}

// getMissingDropboxItems gets the items missing from a dropbox given an array of requested items.
func getMissingDropboxItems(dropbox db.Dropbox, items []db.InventoryItemDropboxOwnerName) error {
	type Body struct {
		Message string   `json:"message"`
		Items   []string `json:"items"`
	}

	idNameArr := make([]string, 0, 0)
	for _, item := range items {
		idNameArr = append(idNameArr, item.InventoryItem.ObjectID, item.InventoryItem.Name)
	}

	jsonBody, err := json.Marshal(&Body{
		Message: "availability",
		Items:   idNameArr,
	})
	if err != nil {
		return err
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	var respBody []byte

	err = retry.Do(func() error {
		resp, err := client.Post(dropbox.URL, "application/json", bytes.NewReader(jsonBody))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		respBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return nil
	}, retry.Attempts(3),
	)
	if err != nil {
		return MissingDropboxError{Dropbox: dropbox}
	}

	if strings.Contains(string(respBody), "cap not found") {
		return MissingDropboxError{Dropbox: dropbox}
	}

	missingItemsRaw := struct {
		MissingItems []string `json:"missing_items"`
	}{}

	err = json.Unmarshal(respBody, &missingItemsRaw)
	if err != nil {
		return err
	}

	missingItems := make([]db.InventoryItemDropboxOwnerName, 0, 0)
	for i := 0; i < len(missingItemsRaw.MissingItems); i += 2 {
		missingItemID := missingItemsRaw.MissingItems[i]
		missingItemName := missingItemsRaw.MissingItems[i+1]
		for _, trueItem := range items {
			if trueItem.InventoryItem.ObjectID == missingItemID && trueItem.InventoryItem.Name == missingItemName {
				missingItems = append(missingItems, trueItem)
			}
		}
	}

	if len(missingItems) == 1 {
		return MissingItemError{Item: missingItems[0]}
	}

	if len(missingItems) > 1 {
		var allErrors []error
		for _, item := range missingItems {
			allErrors = append(allErrors, MissingItemError{Item: item})
		}
		return Dedupe(MissingErrors{Errors: allErrors})
	}

	return nil
}

// commandDropboxSend sends a command to a dropboxURL to send orderItems to targetAvatar.
func commandDropboxSend(dropboxURL string, orderItems []db.OrderItem, targetAvatar string, orderID string) error {
	type Body struct {
		Message      string `json:"message"`
		Items        string `json:"items"`
		TargetAvatar string `json:"target_avatar"`
		OrderID      string `json:"order_id"`
	}

	slDigestibleList := make([]string, 0, 0)
	for _, orderItem := range orderItems {
		slDigestibleList = append(slDigestibleList, orderItem.Name)
	}

	jsonBody, err := json.Marshal(&Body{
		Message:      "send_items",
		Items:        strings.Join(slDigestibleList, ","),
		TargetAvatar: targetAvatar,
		OrderID:      orderID,
	})
	if err != nil {
		return err
	}
	client := http.Client{Timeout: time.Second * 60}
	resp, err := client.Post(dropboxURL, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyString, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(string(bodyString))
	}

	return nil
}
