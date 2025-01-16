package handler

import (
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"shopa/db"
)

func MostLovedRecentBlogs(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	pics, err := db.MostLikedRecentBlogs(c.Request.Context(), dbConn)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "getting most liked recent pictures:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP": c.ClientIP(),
	}).Info("got most liked recent pictures")
	c.JSON(200, pics)
}

func InsertProductPictureLove(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	productPictureID := c.Param("pictureID")
	authAvatar := c.MustGet("authAvatar").(string)
	if productPictureID == "" {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "no productPictureID given")
		return
	}
	err := db.InsertProductPictureLove(c.Request.Context(), dbConn, productPictureID, authAvatar)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "inserting product picture love:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":                 c.ClientIP(),
		"product_picture_id": productPictureID,
	}).Info("inserted product picture love")
	c.JSON(200, nil)
}

func DeleteProductPictureLove(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	productPictureID := c.Param("pictureID")
	authAvatar := c.MustGet("authAvatar").(string)
	if productPictureID == "" {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "no productPictureID given")
		return
	}
	err := db.DeleteProductPictureLove(c.Request.Context(), dbConn, productPictureID, authAvatar)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "deleting product picture love:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":                 c.ClientIP(),
		"product_picture_id": productPictureID,
	}).Info("deleted product picture love")
	c.JSON(200, nil)
}

func InsertBlogPostLove(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	blogPostID := c.Param("postID")
	authAvatar := c.MustGet("authAvatar").(string)
	if blogPostID == "" {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "no blogPostID given")
		return
	}
	err := db.InsertBlogPostLove(c.Request.Context(), dbConn, blogPostID, authAvatar)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "inserting blog post love:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":           c.ClientIP(),
		"blog_post_id": blogPostID,
	}).Info("inserted blog post love")
	c.JSON(200, nil)
}

func DeleteBlogPostLove(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	blogPostID := c.Param("postID")
	authAvatar := c.MustGet("authAvatar").(string)
	if blogPostID == "" {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "no blogPostID given")
		return
	}
	err := db.DeleteBlogPostLove(c.Request.Context(), dbConn, blogPostID, authAvatar)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "deleting blog post love:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":           c.ClientIP(),
		"blog_post_id": blogPostID,
	}).Info("deleted blog post love")
	c.JSON(200, nil)
}
