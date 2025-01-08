package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	ID              string `json:"id"`
	Owner           string `json:"owner"`
	OwnerLegacyName string `json:"owner_legacy_name"`
	Name            string `json:"name" binding:"required"`
	Description     string `json:"description" binding:"required"`

	GroupLink  string `json:"group_link"`
	SLURL      string `json:"sl_url"`
	Flickr     string `json:"flickr"`
	Twitter    string `json:"twitter"`
	ShowAvatar bool   `json:"show_avatar"`
	Banner     string `json:"banner_link" binding:"required"`
	Deleted    bool   `json:"deleted"`
	CSRCSV     string `json:"csr_csv"`
	FacebookLink string `json:"facebook_link"`
}

type StoreProductPreview struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	PrimaryPictureLink string `json:"primary_picture_link"`
}

func StoreOwner(ctx context.Context, conn *pgxpool.Pool, storeID string) (string,error) {
	// get store owner with storeID
	var owner string
	err := conn.QueryRow(ctx, `SELECT owner FROM store WHERE id = $1`, storeID).Scan(&owner)
	if err != nil {
		return "", err
	}
	return owner, nil
}

func DeleteStore(ctx context.Context, conn *pgxpool.Pool, storeID string) error {
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
	DELETE
	FROM customer_order_product_line copl
	WHERE copl.product_id IN (
		SELECT
		p.id
		FROM product p
		join store s2 ON p.store = s2.id
		WHERE s2.id = $1
	)
	`, storeID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
	DELETE
	FROM product_picture pp
	WHERE pp.product_id IN (
		SELECT
		p.id
		FROM product p
		join store s2 ON p.store = s2.id
		WHERE s2.id = $1
	)
	`, storeID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
	DELETE
	FROM product_inventory_line pil
	WHERE pil.product_id IN (
		SELECT
		pil.product_id
		FROM product_inventory_line pil
		join product p ON pil.product_id = p.id
		join store s2 ON p.store = s2.id
		WHERE s2.id = $1
	)
	`, storeID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
	DELETE
	FROM product p
	WHERE p.id IN (
		SELECT
		p.id
		FROM product p
		 JOIN store s2 ON p.store = s2.id
		WHERE s2.id = $1
	)
	`, storeID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
	DELETE FROM store where id = $1
	`, storeID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func GetStoreDetails(ctx context.Context, conn *pgxpool.Pool, storeID string) ([]Store, []StoreProductPreview, error) {
	rows, err := conn.Query(ctx, `
	SELECT store.*, avatar.legacyname FROM store
	JOIN avatar  ON store.owner = avatar.uuid
	WHERE store.id = $1`, storeID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	stores := make([]Store, 0, 0)
	for rows.Next() {
		var store Store
		err = rows.Scan(
			&store.ID,
			&store.Owner,
			&store.Name,
			&store.Description,
			&store.GroupLink,
			&store.SLURL,
			&store.Flickr,
			&store.Twitter,
			&store.ShowAvatar,
			&store.Banner,
			&store.Deleted,
			&store.CSRCSV,
			&store.FacebookLink,
			&store.OwnerLegacyName)
		if err != nil {
			return nil, nil, err
		}
		stores = append(stores, store)
	}

	rows, err = conn.Query(ctx, `
	SELECT
	product.id, product.name, pp.link
	FROM product
	join product_picture pp ON product.id = pp.product_id
	join store s ON product.store = s.id
	where s.id = $1 AND pp.index = 0 AND product.listed = true`, storeID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var prodPreviews []StoreProductPreview
	for rows.Next() {
		var p StoreProductPreview
		err = rows.Scan(&p.ID, &p.Name, &p.PrimaryPictureLink)
		if err != nil {
			return nil, nil, err
		}
		prodPreviews = append(prodPreviews, p)
	}

	return stores, prodPreviews, nil
}

func GetAvatarStores(ctx context.Context, conn *pgxpool.Pool, avatarKey string) ([]Store, error) {
	rows, err := conn.Query(ctx, `
	SELECT store.*, avatar.legacyname FROM store
	JOIN avatar  ON store.owner = avatar.uuid
	WHERE owner = $1`, avatarKey)
	if err != nil {
		return []Store{}, err
	}
	defer rows.Close()

	stores := make([]Store, 0, 0)
	for rows.Next() {
		var store Store
		err = rows.Scan(&store.ID, &store.Owner, &store.Name, &store.Description, &store.GroupLink, &store.SLURL, &store.Flickr, &store.Twitter, &store.ShowAvatar, &store.Banner, &store.Deleted, &store.CSRCSV, &store.FacebookLink, &store.OwnerLegacyName)
		if err != nil {
			return []Store{}, err
		}
		stores = append(stores, store)
	}

	return stores, nil
}

func UpsertStore(ctx context.Context, conn *pgxpool.Pool, store Store) error {
	idReturned := ""
	var err error = nil

	if store.ID == "" {
		err = conn.QueryRow(ctx,
			`INSERT INTO store 
             	VALUES(uuid_generate_v4(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`,
			store.Owner, store.Name, store.Description, store.GroupLink, store.SLURL, store.Flickr, store.Twitter, store.ShowAvatar, store.Banner, store.Deleted, store.CSRCSV, store.FacebookLink,
		).Scan(&idReturned)
	} else {
		_, err = conn.Exec(ctx, `
					UPDATE store SET name = $2, description = $3, group_link = $4, sl_url = $5, flickr = $6, 
						twitter = $7, show_avatar = $8, banner_link = $9, deleted = $10, csr_csv = $11, facebook_link = $12
					WHERE store.id = $1`,
			store.ID,
			store.Name,
			store.Description,
			store.GroupLink,
			store.SLURL,
			store.Flickr,
			store.Twitter,
			store.ShowAvatar,
			store.Banner,
			store.Deleted,
			store.CSRCSV,
			store.FacebookLink,
		)
	}
	if err != nil {
		return err
	}
	if idReturned == "" {
		idReturned = store.ID
	}
	return nil
}
