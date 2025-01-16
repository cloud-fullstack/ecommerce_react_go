package handler

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"io"
	"math"
	"net/http"
	"shopa/db"
	"strings"
	"time"
)

func EditBlogPost(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	authAvatar := c.MustGet("authAvatar").(string)
	in := struct {
		BlogPostID  string `json:"blog_post_id"`
		ContentLink string `json:"content_link"`
		PictureLink string `json:"picture_link"`
		Type        int    `json:"type"`
	}{}
	err := c.ShouldBindJSON(&in)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding blog post:"+err.Error())
		return
	}

	blogPostToEdit, err := db.OneBlogPost(c.Request.Context(), dbConn, in.BlogPostID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "grabbing blogpost:"+err.Error())
		return
	}
	if blogPostToEdit.Owner != authAvatar {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 400, "you are not the owner")
		return
	}

	err = db.EditBlogPost(c.Request.Context(), dbConn, db.BlogPost{
		ID:          in.BlogPostID,
		ContentLink: in.ContentLink,
		PictureLink: in.PictureLink,
		Type:        in.Type,
	})
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "inserting blog post:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":       c.ClientIP(),
		"BlogPost": in,
	}).Info("edited blogpost")
	c.JSON(200, "")
}

// InsertBlogPost inserts a blog post and its picture.
func InsertBlogPost(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	authAvatar := c.MustGet("authAvatar").(string)
	in := struct {
		ProductID   string `json:"product_id"`
		ContentLink string `json:"content_link"`
		PictureLink string `json:"picture_link"`
		Type        int    `json:"type"`
	}{}
	err := c.ShouldBindJSON(&in)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding blog post:"+err.Error())
		return
	}

	if boughtProduct, err := db.BoughtProduct(c.Request.Context(), dbConn, in.ProductID, authAvatar); !boughtProduct {
		if err != nil {
			logRespondError(c, log.Fields{
				"IP":   c.ClientIP(),
				"blog": in,
			}, 400, "determining if buyer bought product:"+err.Error())
			return
		}
		logRespondError(c, log.Fields{
			"IP":   c.ClientIP(),
			"blog": in,
		}, 400, "you have not bought this product")
		return
	}

	blogsByAuthAvatar, err := db.BlogsByAvatarKey(c.Request.Context(), dbConn, authAvatar)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP":   c.ClientIP(),
			"blog": in,
		}, 400, "determining if blogger has already blogged:"+err.Error())
		return
	}
	for _, blog := range blogsByAuthAvatar {
		if blog.ProductID == in.ProductID {
			logRespondError(c, log.Fields{
				"IP":   c.ClientIP(),
				"blog": in,
			}, 400, "blog has already been written for this product")
			return
		}
	}

	id, err := db.InsertBlogPost(c.Request.Context(), dbConn, db.BlogPost{
		Owner:       authAvatar,
		ProductID:   in.ProductID,
		ContentLink: in.ContentLink,
		PictureLink: in.PictureLink,
		Type:        in.Type,
	})
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "inserting blog post:"+err.Error())
		return
	}

	blogPaid, pricePaid, err := db.BlogHasBeenPaid(c.Request.Context(), dbConn, authAvatar, in.ProductID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "inserting blog post:"+err.Error())
		return
	}

	if !blogPaid {
		amt := int(math.Floor(float64(pricePaid) * 0.02))
		if err := payBlogger(c.Request.Context(), dbConn, authAvatar, amt); err != nil {
			logRespondError(c, log.Fields{
				"IP": c.ClientIP(),
			}, 500, "could not pay for blog post:"+err.Error())
			return
		}
		if err := db.SetCustomerPaidForBlog(c.Request.Context(), dbConn, in.ProductID, authAvatar); err != nil {
			logRespondError(c, log.Fields{
				"IP": c.ClientIP(),
			}, 500, "inserting blog post:"+err.Error())
			return
		}
	}

	log.WithFields(log.Fields{
		"IP":       c.ClientIP(),
		"BlogPost": in,
	}).Info("inserted blogpost")
	c.JSON(200, gin.H{"blog_post_id": id})
}

func payBlogger(ctx context.Context, conn db.Conn, bloggerAvatarID string, amount int) error {
	boxes, err := db.GetDropboxRepos(ctx, conn)
	if err != nil {
		return fmt.Errorf("getting dropbox repos: %w", err)
	}
	if len(boxes) == 0 {
		return fmt.Errorf("getting dropbox repos: no dropbox repos exist")
	}
	box := boxes[0]
	raw := fmt.Sprintf(`{"message":"pay_avatar","avatar":"%s","amount":"%d"}`, bloggerAvatarID, amount)
	r := strings.NewReader(raw)
	req, err := http.NewRequest("POST", box.URL, r)
	if err != nil {
		return err
	}
	cl := http.Client{Timeout: 30 * time.Second}
	resp, err := cl.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("could not read error body %w", err)
		}
		return fmt.Errorf("error finding dropbox repo for blog payment: %d | %s", resp.StatusCode, body)
	}

	return nil
}

func DeleteBlogPost(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	authAvatar := c.MustGet("authAvatar").(string)
	id := c.Param("id")
	if id == "" {
		logRespondError(c, log.Fields{
			"IP":         c.ClientIP(),
			"authAvatar": authAvatar,
		}, 500, "no ID provided")
		return
	}

	post, err := db.OneBlogPost(c.Request.Context(), dbConn, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logRespondError(c, log.Fields{
				"IP":         c.ClientIP(),
				"authAvatar": authAvatar,
				"ID":         id,
			}, 400, "blog post ID does not exist")
			return
		}
		logRespondError(c, log.Fields{
			"IP":         c.ClientIP(),
			"authAvatar": authAvatar,
			"ID":         id,
		}, 500, "getting one blog post:"+err.Error())
		return
	}
	if post.Owner != authAvatar {
		logRespondError(c, log.Fields{
			"IP":         c.ClientIP(),
			"authAvatar": authAvatar,
			"ID":         id,
		}, 400, "you do not own this blog post")
		return
	}

	err = db.DeleteBlogPost(c.Request.Context(), dbConn, id)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP":         c.ClientIP(),
			"authAvatar": authAvatar,
			"ID":         id,
		}, 400, "deleting blogpost:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP":         c.ClientIP(),
		"authAvatar": authAvatar,
		"ID":         id,
	}).Info("deleted blog post")
	c.JSON(200, "deleted")
}

func BlogsByAvatarKey(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	avatarKey := c.Param("avatarKey")

	blogs, err := db.BlogsByAvatarKey(c.Request.Context(), dbConn, avatarKey)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "getting blog posts:"+err.Error())
		return
	}

	avatarLegacyName, err := db.AvatarLegacyName(c.Request.Context(), dbConn, avatarKey)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "getting avatar legacy name in getting blog posts:"+err.Error())
		return
	}

	log.WithFields(log.Fields{
		"IP": c.ClientIP(),
	}).Info("got blog posts")
	c.JSON(200, gin.H{
		"blogs":            blogs,
		"avatarLegacyName": avatarLegacyName,
	})
}