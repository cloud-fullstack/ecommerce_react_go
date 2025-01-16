package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Dropbox struct {
	ObjectID string  `json:"object_id"`
	Owner    string  `json:"owner"`
	Region   string  `json:"region"`
	URL      string  `json:"url"`
	PosX     float32 `json:"pos_x"`
	PosY     float32 `json:"pos_y"`
	PosZ     float32 `json:"pos_z"`
}

type DropboxInventoryItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ObjectID    string `json:"object_id"`  /* given by SL's servers */
	DropboxID   string `json:"dropbox_id"` /* for the same object ID can be in multiple dropboxes */
	Perms       int    `json:"perms"`
	Copyable    bool   `json:"copyable"`
	AcquireTime int    `json:"acquire_time"`
}

type DropboxInventoryItemWithDemo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ObjectID    string `json:"object_id"`  /* given by SL's servers */
	DropboxID   string `json:"dropbox_id"` /* for the same object ID can be in multiple dropboxes */
	Perms       int    `json:"perms"`
	Copyable    bool   `json:"copyable"`
	AcquireTime int    `json:"acquire_time"`
	DemoItem    bool   `json:"demo_item"`
}

func GetDropboxRepos(ctx context.Context, conn Conn) ([]Dropbox, error) {
	rows, err := conn.Query(ctx, "SELECT * from dropbox_repo")
	defer rows.Close()
	if err != nil {
		return []Dropbox{}, err
	}

	dropboxes := make([]Dropbox, 0, 0)
	for rows.Next() {
		var dropbox Dropbox
		err = rows.Scan(&dropbox.ObjectID, nil, &dropbox.Owner, &dropbox.Region, &dropbox.URL, &dropbox.PosX, &dropbox.PosY, &dropbox.PosZ)
		if err != nil {
			return []Dropbox{}, err
		}
		dropboxes = append(dropboxes, dropbox)
	}

	return dropboxes, nil
}

func GetAvatarDropboxContents(ctx context.Context, conn *pgxpool.Pool, dropboxUUID string) ([]DropboxInventoryItem, error) {
	rows, err := conn.Query(ctx,
		"SELECT * from inventory_item WHERE dropbox_id = $1",
		dropboxUUID,
	)
	defer rows.Close()
	if err != nil {
		return []DropboxInventoryItem{}, err
	}

	dropboxItems := make([]DropboxInventoryItem, 0, 0)
	for rows.Next() {
		var dropbox DropboxInventoryItem
		var t time.Time
		err = rows.Scan(&dropbox.ID, &dropbox.Name, &dropbox.ObjectID, &dropbox.DropboxID, &dropbox.Perms, &dropbox.Copyable, &t)
		if err != nil {
			return []DropboxInventoryItem{}, err
		}
		dropbox.AcquireTime = int(t.Unix())
		dropboxItems = append(dropboxItems, dropbox)
	}

	return dropboxItems, nil
}

func GetAvatarDropboxes(ctx context.Context, conn *pgxpool.Pool, avatarUUID string) ([]Dropbox, error) {
	rows, err := conn.Query(ctx,
		"SELECT * from dropbox WHERE owner = $1",
		avatarUUID,
	)
	defer rows.Close()
	if err != nil {
		return []Dropbox{}, err
	}

	dropboxes := make([]Dropbox, 0, 0)
	for rows.Next() {
		var dropbox Dropbox
		err = rows.Scan(&dropbox.ObjectID, nil, &dropbox.Owner, &dropbox.Region, &dropbox.URL, &dropbox.PosX, &dropbox.PosY, &dropbox.PosZ)
		if err != nil {
			return []Dropbox{}, err
		}
		dropboxes = append(dropboxes, dropbox)
	}

	return dropboxes, nil
}

func InsertDropbox(ctx context.Context, conn *pgxpool.Pool, info Dropbox) error {
	_, err := conn.Exec(ctx,
		`INSERT INTO dropbox VALUES($1, to_timestamp($2), $3, $4, $5, $6, $7, $8)
			ON CONFLICT ON CONSTRAINT dropbox_pkey
			DO
			UPDATE SET url = $5, pos_x = $6, pos_y = $7, pos_z = $8 
			WHERE dropbox.object_id = $1`,
		info.ObjectID,
		time.Now().Unix(),
		info.Owner,
		info.Region,
		info.URL,
		info.PosX,
		info.PosY,
		info.PosZ)
	if err != nil {
		return err
	}
	return nil
}

func InsertDropboxRepo(ctx context.Context, conn *pgxpool.Pool, info Dropbox) error {
	_, err := conn.Exec(ctx,
		`INSERT INTO dropbox_repo VALUES($1, to_timestamp($2), $3, $4, $5, $6, $7, $8)
			ON CONFLICT ON CONSTRAINT dropbox_repo_pkey
			DO
			UPDATE SET url = $5, pos_x = $6, pos_y = $7, pos_z = $8 
			WHERE dropbox_repo.object_id = $1`,
		info.ObjectID,
		time.Now().Unix(),
		info.Owner,
		info.Region,
		info.URL,
		info.PosX,
		info.PosY,
		info.PosZ)
	if err != nil {
		return err
	}
	return nil
}

// DeleteDropbox deletes the dropbox from the db and any inventory items and product_inventory_lines as a result.
func DeleteDropbox(ctx context.Context, conn Conn, dropboxID string) error {
	_, err := conn.Exec(ctx, `
		DELETE from product_inventory_line AS pil
		WHERE pil.inventory_item_id IN 
		(
			SELECT pil.inventory_item_id
			FROM product_inventory_line AS pil
				 JOIN inventory_item ii ON ii.id = pil.inventory_item_id
				join dropbox d ON ii.dropbox_id = d.object_id
			WHERE d.object_id = $1
		)
		`,
		dropboxID,
	)
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx,
		`DELETE FROM inventory_item AS ii
		WHERE ii.id IN 
		(
			SELECT ii.id
			FROM inventory_item AS ii
				join dropbox d ON ii.dropbox_id = d.object_id
			WHERE d.object_id = $1
		)`,
		dropboxID,
	)
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx,
		`DELETE FROM dropbox WHERE object_id = $1`,
		dropboxID,
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteInventoryItem deletes the inventory item from the db and any product_inventory_lines as a result.
func DeleteInventoryItem(ctx context.Context, conn Conn, inventoryItemID string) error {
	_, err := conn.Exec(ctx,
		`
		DELETE from product_inventory_line AS pil
		WHERE pil.inventory_item_id IN
		(
		  SELECT pil.inventory_item_id
		  FROM product_inventory_line AS pil
				   JOIN inventory_item ii ON ii.id = pil.inventory_item_id
		  WHERE ii.id = $1
		)`,
		inventoryItemID,
	)
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx,
		`DELETE FROM inventory_item WHERE id = $1`,
		inventoryItemID,
	)
	if err != nil {
		return err
	}
	return nil
}

// UpdateDropboxContents and returns orphaned inventory items and products relating to those inventory items.
func UpdateDropboxContents(ctx context.Context, conn *pgxpool.Pool, items []DropboxInventoryItem) ([]DropboxInventoryItem, []Product, error) {
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return nil, nil, err
	}

	for _, item := range items {
		_, err = tx.Exec(ctx,
			`INSERT INTO inventory_item VALUES(uuid_generate_v4(), $1, $2, $3, $4, $5, to_timestamp($6))
				 ON CONFLICT ON CONSTRAINT unique_pair
				 DO
				 UPDATE SET name = $1, perms = $4, copyable = $5 
				 WHERE inventory_item.object_id = $2 AND inventory_item.dropbox_id = $3 AND inventory_item.acquire_time = to_timestamp($6)`,
			item.Name,
			item.ObjectID,
			item.DropboxID,
			item.Perms,
			item.Copyable,
			item.AcquireTime,
		)
		if err != nil {
			return nil, nil, err
		}
	}

	// inclues copyableOrphanRecords and noCopyOrphanRecords
	var allOrphanRecords []DropboxInventoryItem
	var allProductsWithOrphans []Product

	if len(items) > 0 {
		dropboxID := items[0].DropboxID
		rows, err := tx.Query(ctx, `SELECT * FROM inventory_item
					   WHERE dropbox_id = $1`, dropboxID)
		if err != nil {
			return nil, nil, err
		}
		itemsInDB := make(map[string]DropboxInventoryItem)
		for rows.Next() {
			itemInDB := DropboxInventoryItem{}
			var acquireTime time.Time
			err = rows.Scan(&itemInDB.ID, &itemInDB.Name, &itemInDB.ObjectID, &itemInDB.DropboxID, &itemInDB.Perms, &itemInDB.Copyable, &acquireTime)
			if err != nil {
				return nil, nil, err
			}
			itemInDB.AcquireTime = int(acquireTime.Unix())
			itemsInDB[itemInDB.ID] = itemInDB
		}

		copyableOrphanRecords := make(map[string]DropboxInventoryItem)
		for _, itemInDropbox := range items {
			deleteItem := false
			deleteKey := ""
			for dbID, itemInDB := range itemsInDB {
				if itemInDB.ObjectID == itemInDropbox.ObjectID && itemInDB.AcquireTime == itemInDropbox.AcquireTime {
					deleteItem = true
					deleteKey = dbID
					continue
				}
			}
			if deleteItem {
				if !itemsInDB[deleteKey].Copyable {
					// find the correct one by acquire time
				}
				delete(itemsInDB, deleteKey)
			}
		}
		copyableOrphanRecords = itemsInDB

		noCopyOrphanRecords := make([]DropboxInventoryItem, 0, len(copyableOrphanRecords))
		for key, item := range copyableOrphanRecords {
			if !item.Copyable {
				noCopyOrphanRecords = append(noCopyOrphanRecords, item)
				delete(copyableOrphanRecords, key)
			}
		}

		for _, record := range copyableOrphanRecords {
			allOrphanRecords = append(allOrphanRecords, record)
		}
		for _, record := range noCopyOrphanRecords {
			allOrphanRecords = append(allOrphanRecords, record)
		}

		for _, item := range allOrphanRecords {
			p, err := ProductsAssociatedWithInventoryItem(ctx, conn, item.ID)
			if err != nil {
				return nil, nil, err
			}
			allProductsWithOrphans = append(allProductsWithOrphans, p...)
		}

		if len(noCopyOrphanRecords) > 0 {
			for _, record := range noCopyOrphanRecords {
				_, err = tx.Exec(ctx, `DELETE from product_inventory_line AS pil
					WHERE pil.inventory_item_id IN 
					(
						SELECT pil.inventory_item_id
						FROM product_inventory_line AS pil
							 JOIN inventory_item ii ON ii.id = pil.inventory_item_id
						WHERE ii.object_id = $1 AND ii.name = $2
					)`,
					record.ObjectID, record.Name)
				if err != nil {
					return nil, nil, err
				}
			}
		}

		if len(copyableOrphanRecords) > 0 {
			for _, record := range copyableOrphanRecords {
				_, err = tx.Exec(ctx, `DELETE from product_inventory_line AS pil
					WHERE pil.inventory_item_id IN 
					(
						SELECT pil.inventory_item_id
						FROM product_inventory_line AS pil
							 JOIN inventory_item ii ON ii.id = pil.inventory_item_id
						WHERE ii.object_id = $1 AND ii.name = $2
					)`,
					record.ObjectID, record.Name)
				if err != nil {
					return nil, nil, err
				}
			}
		}

		if len(copyableOrphanRecords) > 0 {
			for _, record := range copyableOrphanRecords {
				_, err = tx.Exec(ctx, `DELETE FROM inventory_item WHERE object_id = $1 AND name = $2`,
					record.ObjectID, record.Name)
				if err != nil {
					return nil, nil, err
				}
			}
		}

		if len(noCopyOrphanRecords) > 0 {
			for _, record := range noCopyOrphanRecords {
				_, err = tx.Exec(ctx, `DELETE FROM inventory_item WHERE acquire_time = to_timestamp($1) AND name = $2 AND object_id = $3`,
					record.AcquireTime, record.Name, record.ObjectID)
				if err != nil {
					return nil, nil, err
				}
			}
		}

	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, nil, err
	}
	return allOrphanRecords, allProductsWithOrphans, nil
}
