package handler

import (
	"crypto/sha1"
	"fmt"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/minio/highwayhash"
	"strings"
)

// GenToken Generates a token used to authorize a user in a web browser.
// A hash is given by user which is decoded to verify the user is who they say they are.
// todo: update the hash to be given as a header, so logging is easier
func GenToken(c *gin.Context) {
	userInfo := struct {
		Hash       string `json:"hash"`
		LegacyName string `json:"legacy_name"`
		AvatarKey  string `json:"avatar_key"`
	}{}
	err := c.ShouldBindJSON(&userInfo)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding userInfo:"+err.Error())
		return
	}

	// todo: make salt specific to user
	salt := "shopaTMAC#3"
	h := sha1.New()
	h.Write([]byte(userInfo.LegacyName + "." + salt))
	expectedHash := fmt.Sprintf("%x", string(h.Sum(nil)))

	// todo: refactor into validation process
	if userInfo.Hash != expectedHash {
		logRespondError(c, log.Fields{
			"IP":         c.ClientIP(),
			"avatarKey":  userInfo.AvatarKey,
			"legacyName": userInfo.LegacyName,
		}, 500, "hash given is unexpected:"+err.Error())
		return
	}

	authToken := MakeSessionKey(userInfo.AvatarKey)
	log.WithFields(log.Fields{
		"IP":         c.ClientIP(),
		"avatarKey":  userInfo.AvatarKey,
		"legacyName": userInfo.LegacyName,
	}).Info("gave user auth token")
	c.JSON(200, gin.H{
		"token": authToken,
		"uuid":  userInfo.AvatarKey,
	})
}

// MiddlewareAuth passes onto the next handler if the auth provided (as "hash.avatarKey") is valid.
func MiddlewareAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || strings.Index(authHeader, ".") == -1 {
			log.WithFields(log.Fields{"IP": c.ClientIP(), "url": c.Request.URL.String()}).Errorf("failed authentication. no Authorization header")
			c.JSON(200, gin.H{
				"error":   true,
				"message": "failed authentication",
			})
			c.Abort()
			return
		}
		split := strings.Split(authHeader, ".")
		hash := split[0]
		aviUUID := split[1]
		if !ValidSessionKey(hash, aviUUID) {
			log.WithFields(log.Fields{"IP": c.ClientIP(), "uuid": aviUUID, "url": c.Request.URL.String()}).Errorf("failed authentication")
			c.JSON(200, gin.H{
				"error":   true,
				"message": "failed authentication",
			})
			c.Abort()
			return
		}
		c.Set("authAvatar", aviUUID)
		c.Next()
	}
}

// MakeSessionKey accepts an avatar's key/UUID (with/without dashes) to make a web browser auth key.
func MakeSessionKey(avatarKey string) string {
	avatarKey = strings.ReplaceAll(avatarKey, "-", "")
	hasher, err := highwayhash.New([]byte("12121212121212121212121212121212"))
	if err != nil {
		return ""
	}

	hasher.Write([]byte(avatarKey + "shopa1756"))
	hash := fmt.Sprintf("%x", string(hasher.Sum(nil)))
	return hash
}

func ValidSessionKey(hashB64 string, avatarKey string) bool {
	hash := MakeSessionKey(avatarKey)
	if hashB64 == hash {
		return true
	}
	return false
}
