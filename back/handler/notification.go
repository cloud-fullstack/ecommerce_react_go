package handler

import (
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"shopa/db"
)

func GetNotification(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	authAvatar := c.MustGet("authAvatar").(string)
	notifications, err := db.GetNotifications(c.Request.Context(), dbConn, authAvatar)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "getting notifications:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":            c.ClientIP(),
		"avatar":        authAvatar,
		"notifications": notifications,
	}).Info("requested notifications")
	c.JSON(200, notifications)
}
