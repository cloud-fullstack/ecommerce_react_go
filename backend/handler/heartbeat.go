package handler

import (
	"context"
	"fmt"
	"github.com/apex/log"
	"github.com/avast/retry-go"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
	"shopa/db"
	"strings"
	"time"
)

// pingHUD pings the user's HUD in-world to determine if it is alive.
func pingHUD(ctx context.Context, conn *pgxpool.Pool, avatarKey string) (bool, error) {
	row := conn.QueryRow(ctx, "SELECT hud_url FROM avatar WHERE uuid=$1", avatarKey)
	hudURL := ""
	err := row.Scan(&hudURL)
	if err != nil {
		return false, err
	}

	var statusCode int
	var status string

	err = retry.Do(func() error {
		client := http.Client{
			Timeout: 5 * time.Second,
		}
		resp, err := client.Get(hudURL)
		if err != nil {
			return err
		}
		statusCode = resp.StatusCode
		status = resp.Status
		return nil
	}, retry.Attempts(3),
	)
	if err != nil {
		return false, err
	}

	if statusCode != 200 {
		return false, fmt.Errorf("non-200 status code: " + status)
	}
	return true, nil
}

// PingHUD pings the user's HUD in-world to determine if it is alive.
// It wraps pingHUD in a handler and prevents people from heartbeating HUDs that are not their own.
func PingHUD(c *gin.Context) {
	conn := c.MustGet("dbConn").(*pgxpool.Pool)

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

	authHeader := c.GetHeader("Authorization")
	split := strings.Split(authHeader, ".")
	authAvatarKey := split[1]

	if userInfo.AvatarKey != authAvatarKey {
		// todo: refactor into validation
		logRespondError(c, log.Fields{
			"IP":                   c.ClientIP(),
			"requested_avatar_key": userInfo.AvatarKey,
			"auth_avatar_key":      authAvatarKey,
		}, 500, "cannot heartbeat HUDs that do not belong to you")
		return
	}

	_, err = pingHUD(c.Request.Context(), conn, userInfo.AvatarKey)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP":                   c.ClientIP(),
			"requested_avatar_key": userInfo.AvatarKey,
			"auth_avatar_key":      authAvatarKey,
		}, 500, "failed to heartbeat hud: "+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":              c.ClientIP(),
		"userInfo":        userInfo,
		"auth_avatar_key": authAvatarKey,
	}).Info("pinged HUD")
	c.JSON(200, gin.H{
		"message": "pinged HUD",
	})
}

// UpdateHeartbeatURLHUD updates the URL for an avatar's HUD in the database.
func UpdateHeartbeatURLHUD(c *gin.Context) {
	conn := c.MustGet("dbConn").(*pgxpool.Pool)

	hudInfo := struct {
		URL        string `json:"url"`
		AvatarKey  string `json:"avatar_key"`
		LegacyName string `json:"legacy_name"`
		Region     string `json:"region"`
		Position   string `json:"position"`
	}{}
	err := c.ShouldBindJSON(&hudInfo)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding SLDropbox:"+err.Error())
		return
	}

	// todo: disallow people updating other people's huds

	// todo: change to db.GetAvatar or similar and check for GetAvatar having no avatar
	// must ensure avatar is inserted
	err = db.InsertAvatar(c.Request.Context(), conn, hudInfo.AvatarKey, hudInfo.LegacyName)
	if err != nil && !strings.Contains(err.Error(), "duplicate key") {
		logRespondError(c, log.Fields{
			"IP":      c.ClientIP(),
			"hudInfo": hudInfo,
		}, 500, "inserting avatar: "+err.Error())
		return
	}

	err = db.UpdateHUDURL(c.Request.Context(), conn, hudInfo.URL, hudInfo.AvatarKey)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP":      c.ClientIP(),
			"hudInfo": hudInfo,
		}, 500, "updating hud URL: "+err.Error())
		return
	}

	_, err = pingHUD(c.Request.Context(), conn, hudInfo.AvatarKey)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP":      c.ClientIP(),
			"hudInfo": hudInfo,
		}, 500, "failed to heartbeat HUD: "+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":      c.ClientIP(),
		"hudInfo": hudInfo,
	}).Info("updated heartbeat url and tested heartbeat")
	c.JSON(200, gin.H{
		"message": "updated heartbeat url and tested heartbeat",
	})
}
