package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

type BlogPost struct {
	ID              string `json:"id"`
	Owner           string `json:"owner"`
	ProductID       string `json:"product_id"`
	ContentLink     string `json:"content_link"`
	PictureLink     string `json:"picture_link"`
	Type            int    `json:"type"`
	OwnerLegacyName string `json:"owner_legacy_name"`
	LoveCount       int    `json:"love_count"`
	ProductStoreID  string `json:"product_store_id"`
}

type MostLikedBlog struct {
	PictureID    string `json:"picture_id"`
	PictureLink  string `json:"picture_link"`
	StoreID      string `json:"store_id"`
	ProductID    string `json:"product_id"`
	BlogPostID   string `json:"blog_post_id"`
	ProductName  string `json:"product_name"`
	AuthorName   string `json:"author_name"`
	AuthorId     string `json:"author_id"`
	ExternalLink string `json:"external_link"`
	LoveCount    int    `json:"love_count"`
}

func MostLikedRecentBlogs(ctx context.Context, conn Conn) ([]MostLikedBlog, error) {
	rows, err := conn.Query(ctx, `
	SELECT blog_post_id,
       picture_link AS picture_link,
       s.id AS store_id,
       p.id AS product_id,
	   pp.id,		
       p.name AS product_name,
	   a.legacyname as author_name,
	   a.uuid as author_id,
	   pp.content_link as external_link,
       COUNT(blog_post_love.blog_post_id) AS love_count
	FROM blog_post_love
			 JOIN blog_post pp ON blog_post_love.blog_post_id = pp.id
			 JOIN product p ON pp.product_id = p.id
			 JOIN store s ON p.store = s.id
			 JOIN avatar a on pp.owner = a.uuid
	GROUP BY blog_post_id, picture_link, s.id, p.id, p.name, pp.id, a.legacyname, a.uuid, external_link
	ORDER BY love_count DESC
	LIMIT 9`)
	// old version with time
	// SELECT blog_post_id,
	//    picture_link AS picture_link,
	//    s.id AS store_id,
	//    p.id AS product_id,
	//    pp.id,
	//    p.name AS product_name,
	//    a.legacyname as author_name,
	//    a.uuid as author_id,
	//    pp.content_link as external_link,
	//    COUNT(blog_post_love.blog_post_id) AS love_count
	// FROM blog_post_love
	// 		 JOIN blog_post pp ON blog_post_love.blog_post_id = pp.id
	// 		 JOIN product p ON pp.product_id = p.id
	// 		 JOIN store s ON p.store = s.id
	// 		 JOIN avatar a on pp.owner = a.uuid
	// WHERE blog_post_love.date >= NOW() - INTERVAL '72 HOURS'
	// GROUP BY blog_post_id, picture_link, s.id, p.id, p.name, pp.id, a.legacyname, a.uuid, external_link
	// ORDER BY love_count DESC
	// LIMIT 5
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	loves := []MostLikedBlog{}
	for rows.Next() {
		var love MostLikedBlog
		err = rows.Scan(
			&love.PictureID,
			&love.PictureLink,
			&love.StoreID,
			&love.ProductID,
			&love.BlogPostID,
			&love.ProductName,
			&love.AuthorName,
			&love.AuthorId,
			&love.ExternalLink,
			&love.LoveCount,
		)
		if err != nil {
			return nil, err
		}
		loves = append(loves, love)
	}
	return loves, nil
}

func BlogsByAvatarKey(ctx context.Context, conn Conn, avatarID string) ([]BlogPost, error) {
	rows, err := conn.Query(ctx, `
	SELECT
		blog_post.*,
		legacyname as blog_post_owner,
		COUNT(blog_post_love) as love_count,
		p.store as product_store_id
	FROM blog_post
		JOIN product p ON p.id = blog_post.product_id
		JOIN avatar ON blog_post.owner = avatar.uuid
		LEFT JOIN blog_post_love on blog_post.id = blog_post_love.blog_post_id
	WHERE blog_post.owner = $1
	GROUP BY blog_post.id,
		blog_post.owner,
		blog_post.product_id,
		blog_post.content_link,
		blog_post.picture_link,
		blog_post.type,
		legacyname, p.store`, avatarID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var posts []BlogPost
	for rows.Next() {
		var post BlogPost
		if err := rows.Scan(
			&post.ID,
			&post.Owner,
			&post.ProductID,
			&post.ContentLink,
			&post.PictureLink,
			&post.Type,
			&post.OwnerLegacyName,
			&post.LoveCount,
			&post.ProductStoreID,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// ProductBlogPosts returns the blog posts associated with productID.
func ProductBlogPosts(ctx context.Context, conn Conn, productID string) ([]BlogPost, error) {
	rows, err := conn.Query(ctx, `
	SELECT
		blog_post.*,
		legacyname as blog_post_owner,
		COUNT(blog_post_love) as love_count,
		p.store as product_store_id
	FROM blog_post
		JOIN product p ON p.id = blog_post.product_id
		JOIN avatar ON blog_post.owner = avatar.uuid
		LEFT JOIN blog_post_love on blog_post.id = blog_post_love.blog_post_id
	WHERE p.id = $1
	GROUP BY blog_post.id,
		blog_post.owner,
		blog_post.product_id,
		blog_post.content_link,
		blog_post.picture_link,
		blog_post.type,
		legacyname, p.store`, productID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var posts []BlogPost
	for rows.Next() {
		var post BlogPost
		if err := rows.Scan(
			&post.ID,
			&post.Owner,
			&post.ProductID,
			&post.ContentLink,
			&post.PictureLink,
			&post.Type,
			&post.OwnerLegacyName,
			&post.LoveCount,
			&post.ProductStoreID,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func BlogPostIDsLovedByAvatar(ctx context.Context, conn Conn, avatarViewing string) ([]string, error) {
	if avatarViewing == "" {
		return nil, nil
	}

	rows, err := conn.Query(ctx, `SELECT blog_post_id FROM blog_post_love WHERE owner = $1`, avatarViewing)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// InsertBlogPost inserts the blog post and returns the new posts ID.
func InsertBlogPost(ctx context.Context, conn Conn, post BlogPost) (string, error) {
	var id string
	err := conn.QueryRow(ctx, `
	INSERT INTO blog_post 
	VALUES(uuid_generate_v4(), $1, $2, $3, $4, $5) RETURNING id`,
		post.Owner,
		post.ProductID,
		post.ContentLink,
		post.PictureLink,
		post.Type,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// BlogHasBeenPaid returns if the blog was paid, and the amount paid for the product.
func BlogHasBeenPaid(ctx context.Context, conn Conn, blogWriterAvatarID, productID string) (bool, int, error) {
	query := fmt.Sprintf(`select paid_review, receipt -> 0 -> 'price'  from customer_order where buyer = '%s' AND (receipt -> 0 -> 'product_id')::text = '"%s"'`, blogWriterAvatarID, productID)
	rows, err := conn.Query(ctx, query)
	if err != nil {
		return false, 0, err
	}
	defer rows.Close()
	var price int
	for rows.Next() {
		var paid sql.NullBool
		if err := rows.Scan(&paid, &price); err != nil {
			return false, 0, err
		}
		if paid.Bool {
			return true, price, nil
		}
	}
	return false, price, nil
}

// SetCustomerPaidForBlog sets a customer_order for a blog post to paid, so they can't get paid for it again.
func SetCustomerPaidForBlog(ctx context.Context, conn Conn, productID, bloggerAvatarID string) error {
	query := fmt.Sprintf(`UPDATE customer_order set paid_review = true where buyer = '%s' AND (receipt -> 0 -> 'product_id')::text = '"%s"'`, bloggerAvatarID, productID)
	_, err := conn.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

// EditBlogPost edits the blog post.
func EditBlogPost(ctx context.Context, conn Conn, post BlogPost) error {
	_, err := conn.Exec(ctx, `
	UPDATE blog_post SET content_link = $1, picture_link = $2, type = $3`,
		post.ContentLink,
		post.PictureLink,
		post.Type,
	)
	if err != nil {
		return err
	}
	return nil
}

// DeleteBlogPost deletes the blog post and its associated loves.
func DeleteBlogPost(ctx context.Context, conn Conn, id string) error {
	_, err := conn.Exec(ctx, `DELETE FROM blog_post_love WHERE blog_post_id = $1`, id)
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, `DELETE FROM blog_post WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

// OneBlogPost returns the blog post associated with ID.
func OneBlogPost(ctx context.Context, conn Conn, id string) (BlogPost, error) {
	var post BlogPost
	err := conn.QueryRow(ctx, `SELECT * FROM blog_post WHERE id = $1`, id).Scan(
		&post.ID,
		&post.Owner,
		&post.ProductID,
		&post.ContentLink,
		&post.PictureLink,
		&post.Type,
	)
	if err != nil {
		return BlogPost{}, err
	}
	return post, nil
}

func InsertBlogPostLove(ctx context.Context, conn Conn, postID, avatarID string) error {
	loved, err := BlogPostLovedByAvatar(ctx, conn, postID, avatarID)
	if err != nil {
		return err
	}
	if loved {
		return errors.New("already loved by avatar")
	}
	_, err = conn.Exec(ctx, `INSERT INTO blog_post_love 
		VALUES(uuid_generate_v4(), $1, $2, to_timestamp($3))`,
		avatarID, postID, time.Now().Unix(),
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteBlogPostLove(ctx context.Context, conn Conn, postID, avatarID string) error {
	loved, err := BlogPostLovedByAvatar(ctx, conn, postID, avatarID)
	if err != nil {
		return err
	}
	if !loved {
		return errors.New("post is not loved by avatar")
	}

	_, err = conn.Exec(ctx, `DELETE FROM blog_post_love 
		WHERE $1 = blog_post_id AND owner = $2`, postID, avatarID,
	)
	if err != nil {
		return err
	}
	return nil
}

func BlogPostLovedByAvatar(ctx context.Context, conn Conn, postID string, avatarID string) (bool, error) {
	var loved bool
	err := conn.QueryRow(ctx, `SELECT owner = $2 from 
		blog_post_love 
		WHERE blog_post_id = $1 AND owner = $2`, postID, avatarID).Scan(&loved)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
	}
	return loved, nil
}
