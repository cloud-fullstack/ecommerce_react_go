package handler

import (
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4/pgxpool"
	"shopa/db"
)

// StoreDetails gets a store's details and its products
func StoreDetails(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	storeID := c.Param("storeID")

	storeDetails, prodPreviews, err := db.GetStoreDetails(c.Request.Context(), dbConn, storeID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding store:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":               c.ClientIP(),
		"store_id":         storeID,
		"store_details":    storeDetails,
		"product_previews": prodPreviews,
	}).Info("got store previews")
	c.JSON(200, gin.H{
		"store_details":    storeDetails,
		"product_previews": prodPreviews,
	})
}

func GetAvatarStores(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	userInfo := struct {
		AvatarKey string `json:"avatar_key"`
	}{}
	err := c.ShouldBindJSON(&userInfo)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding store:"+err.Error())
		return
	}

	stores, err := db.GetAvatarStores(c.Request.Context(), dbConn, userInfo.AvatarKey)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding store:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":     c.ClientIP(),
		"stores": stores,
	}).Info("got avatar stores")
	c.JSON(200, stores)
}

func InsertStore(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	authAvatar := c.MustGet("authAvatar").(string)
	store := db.Store{}
	err := c.ShouldBindJSON(&store)
	if err != nil {
		if fe, ok := err.(validator.ValidationErrors); ok {
			err1 := fe[0]
			logRespondError(c, log.Fields{
				"IP": c.ClientIP(),
			}, 500, translate(err1), gin.H{
				"field":err1.Field(),
			})
			return
		}

		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding store:"+err.Error())
		return
	}

	storeOwner := ""
	if store.ID != "" {
		storeOwner, err = db.StoreOwner(c.Request.Context(), dbConn, store.ID)
		if err != nil {
			logRespondError(c, log.Fields{
				"IP": c.ClientIP(),
			}, 500, "determining store owner: " + err.Error())
			return
		}
	}

	if authAvatar != storeOwner && store.ID != "" {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "cannot edit store you do not own")
		return
	}

	if store.Name == "" {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "store name is empty")
		return
	}

	err = db.UpsertStore(c.Request.Context(), dbConn, store)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "upserting store:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":    c.ClientIP(),
		"store": store,
	}).Info("inserted/updated store")
	c.JSON(200, gin.H{
		"message": "inserted/updated store",
	})
}

func DeleteStore(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	authAvatar := c.MustGet("authAvatar").(string)
	storeID := c.Param("storeID")

	if storeID == "" {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "storeID not present")
		return
	}

	storeOwner, err := db.StoreOwner(c.Request.Context(), dbConn, storeID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "determining details owner: " + err.Error())
		return
	}

	if storeOwner != authAvatar {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "storeOwner mismatch")
		return
	}

	err = db.DeleteStore(c.Request.Context(), dbConn, storeID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "deleting store from db:" + err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":    c.ClientIP(),
		"store": storeID,
		"owner": authAvatar,
	}).Info("deleted store")
	c.JSON(200, gin.H{
		"message": "deleted store",
	})
}