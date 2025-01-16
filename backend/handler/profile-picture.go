package handler

import (
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

func GrabProfilePicture(c *gin.Context) {
	requestingAvatarKey := strings.Split(c.Request.Header.Get("Authorization"), ".")[1]

	userInfo := struct {
		AvatarKey string `json:"avatar_key"`
	}{}
	err := c.ShouldBindJSON(&userInfo)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding userInfo: "+err.Error())
		return
	}

	resp, err := http.Get("http://world.secondlife.com/resident/" + userInfo.AvatarKey)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP":                  c.ClientIP(),
			"requestingAvatarKey": requestingAvatarKey,
		}, 500, "reading webpage when grabbing profile picture: "+err.Error())
		return
	}

	html, _ := ioutil.ReadAll(resp.Body)

	// grabs the url inside
	// <img alt="profile image" src="https://secondlife.com/app/image/5295f42d-4871-c61c-ddd7-d57d44fa0494/1" class="parcelimg">
	urlStart := strings.Index(string(html), `<img alt="profile image" src="`)
	urlEnd := strings.Index(string(html), `" class="parcelimg"`)

	if urlStart == -1 || urlEnd == -1 {
		logRespondError(c, log.Fields{
			"IP":                  c.ClientIP(),
			"requestingAvatarKey": requestingAvatarKey,
		}, 500, "did not find picture in page")
		return
	}

	preRedirectPicURL := string(html)[urlStart+len(`<img alt="profile image" src="`) : urlEnd]

	noRedirectClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	resp, err = noRedirectClient.Get(preRedirectPicURL)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP":                  c.ClientIP(),
			"requestingAvatarKey": requestingAvatarKey,
		}, 500, "doing an http request to get true picture:"+err.Error())
		return
	}

	// because LL redirects, and we only want to give user the URL to the image, not the image itself
	pic := resp.Header.Get("location")

	log.WithFields(log.Fields{
		"IP": c.ClientIP(), "requestingAvatarKey": requestingAvatarKey,
		"targetAvatarKey": userInfo.AvatarKey,
	}).Info("gave user profile picture")
	c.JSON(200, gin.H{
		"profile_picture": pic,
	})
}
