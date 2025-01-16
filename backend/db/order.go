package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type OrderLine struct {
	ProductID string `json:"product_id"`
	Demo      bool   `json:"demo"`
}

type InvoiceLine struct {
	SellerAvatarUUID string `json:"seller_avatar_uuid"`
	SellerAvatarName string `json:"seller_avatar_name"`
	ProductsSum      int    `json:"products_sum"`
}

type OrderItem struct {
	ObjectID   string `json:"object_id"`
	Name       string `json:"name"`
	DropboxURL string `json:"dropbox_url"`
	DatabaseID string `json:"database_id"`
	Copyable   bool   `json:"copyable"`
}

type ReceiptLine struct {
	OrderID                string `json:"order_id"`
	BuyerID                string `json:"buyer_id"`
	ProductID              string `json:"product_id"`
	ProductOwnerID         string `json:"product_owner_id"`
	ProductOwnerLegacyName string `json:"product_owner_legacy_name"`
	ProductName            string `json:"product_name"`
	StoreID                string `json:"store_id"`
	StoreName              string `json:"store_name"`
	Demo                   bool   `json:"demo"`
	DiscountActive         bool   `json:"discount_active"`
	DiscountedPrice        int    `json:"discounted_price"`
	Price                  int    `json:"price"`
	TruePrice              int    `json:"true_price"`
}

func CreateOrder(ctx context.Context, conn *pgxpool.Pool, orderLines []OrderLine, avatarTargetID string) (orderID string, err error) {
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return "", err
	}

	var returnedID string
	// todo: change uuid to shortuuid or a snowflake
	err = tx.QueryRow(ctx, `INSERT into customer_order VALUES(uuid_generate_v4(),$1, $2, $3, $4, $5) RETURNING id`,
		avatarTargetID, time.Now(), false, false, false).Scan(&returnedID)
	if err != nil {
		return "", err
	}

	for _, line := range orderLines {
		_, err = tx.Exec(ctx, `INSERT into customer_order_product_line VALUES($1, $2, $3)`, returnedID, line.ProductID, line.Demo)
		if err != nil {
			return "", err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", err
	}

	receiptLines, err := OrderReceipt(ctx, conn, returnedID)
	if err != nil {
		return "", err
	}

	err = SetOrderReceiptJSON(ctx, conn, returnedID, receiptLines)
	if err != nil {
		return "", err
	}
	return returnedID, nil
}

// OrderPayables gets the amount owed to each merchant involved in an order.
func OrderPayables(ctx context.Context, conn *pgxpool.Pool, orderID string) ([]InvoiceLine, error) {
	rows, err := conn.Query(ctx, `
		SELECT a.uuid,
			   a.legacyname,
			   SUM(CASE
				   WHEN demo THEN 0
				   WHEN discount_active THEN discounted_price
				   ELSE price
			   END) AS  products_sum
		FROM customer_order
		JOIN customer_order_product_line copl ON customer_order.id = copl.customer_order
		JOIN product p ON copl.product_id = p.id
		JOIN avatar a ON p.owner = a.uuid
		WHERE customer_order.id = $1
		GROUP BY a.uuid, a.legacyname`,
		orderID)
	if err != nil {
		return nil, err
	}
	var invoiceLines []InvoiceLine
	for rows.Next() {
		var invoiceLine InvoiceLine
		err = rows.Scan(&invoiceLine.SellerAvatarUUID, &invoiceLine.SellerAvatarName, &invoiceLine.ProductsSum)
		if err != nil {
			return nil, err
		}
		invoiceLines = append(invoiceLines, invoiceLine)
	}
	return invoiceLines, nil
}

func MarkOrderFulfilled(ctx context.Context, conn Conn, orderID string) error {
	_, err := conn.Exec(ctx, `UPDATE customer_order SET fulfilled = true WHERE id = $1`, orderID)
	if err != nil {
		return err
	}
	return nil
}

func SetOrderReceiptJSON(ctx context.Context, conn *pgxpool.Pool, orderID string, receiptLines []ReceiptLine) error {
	_, err := conn.Exec(ctx, `UPDATE customer_order SET receipt = $1 WHERE id = $2`, receiptLines, orderID)
	if err != nil {
		return err
	}
	return nil
}

func OrderReceipt(ctx context.Context, conn Conn, orderID string) ([]ReceiptLine, error) {
	rows, err := conn.Query(ctx, `SELECT 
       customer_order.id,
	   customer_order.buyer,
       copl.product_id as product_id,
       p.owner as product_owner,
       a.legacyname as product_owner_legacy_name,
       p.name as product_name,
       s.id as store_id,
       s.name as store_name,
       copl.demo,
       p.discount_active,
       p.discounted_price,
       p.price,
       (CASE
            WHEN demo THEN 0
            WHEN discount_active THEN discounted_price
            ELSE price
           END) AS  true_price
		FROM customer_order
		JOIN customer_order_product_line copl ON customer_order.id = copl.customer_order
		JOIN product p ON copl.product_id = p.id
		JOIN avatar a ON p.owner = a.uuid
		JOIN store s ON p.store = s.id
		WHERE customer_order.id = $1`, orderID)
	if err != nil {
		return nil, err
	}
	receipt := make([]ReceiptLine, 0, 0)
	for rows.Next() {
		var rLine ReceiptLine
		err = rows.Scan(
			&rLine.OrderID,
			&rLine.BuyerID,
			&rLine.ProductID,
			&rLine.ProductOwnerID,
			&rLine.ProductOwnerLegacyName,
			&rLine.ProductName,
			&rLine.StoreID,
			&rLine.StoreName,
			&rLine.Demo,
			&rLine.DiscountActive,
			&rLine.DiscountedPrice,
			&rLine.Price,
			&rLine.TruePrice,
		)
		if err != nil {
			return nil, err
		}
		receipt = append(receipt, rLine)
	}
	return receipt, nil
}

type Order struct {
	ID          string        `json:"id"`
	Buyer       string        `json:"buyer"`
	Created     time.Time     `json:"created"`
	Fulfilled   bool          `json:"fulfilled"`
	Cancelled   bool          `json:"cancelled"`
	Paid_review bool          `json:"paid_review"`
	Receipt     []ReceiptLine `json:"receipt"`
}

func OrderHistory(ctx context.Context, conn Conn, avatarKey string) ([]Order, error) {
	rows, err := conn.Query(ctx, `SELECT * FROM customer_order WHERE customer_order.buyer = $1`, avatarKey)
	if err != nil {
		return nil, err
	}
	orders := make([]Order, 0, 0)
	for rows.Next() {
		var order Order
		var receiptLines []ReceiptLine
		err = rows.Scan(
			&order.ID,
			&order.Buyer,
			&order.Created,
			&order.Fulfilled,
			&order.Cancelled,
			&order.Paid_review,
			&receiptLines,
		)
		order.Receipt = receiptLines
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func CancelOrder(ctx context.Context, conn *pgxpool.Pool, orderID string) error {
	_, err := conn.Exec(ctx, `UPDATE customer_order SET cancelled = true WHERE id = $1`, orderID)
	if err != nil {
		return err
	}
	return nil
}

// OrderItems gets the items pertaining to an order.
func OrderItems(ctx context.Context, conn *pgxpool.Pool, orderID string) ([]OrderItem, error) {
	rows, err := conn.Query(ctx, `
		SELECT ii.object_id, ii.name, d.url, ii.id, ii.copyable as database_id
		FROM customer_order
         JOIN customer_order_product_line copl ON customer_order.id = copl.customer_order
         JOIN product p ON copl.product_id = p.id
         JOIN product_inventory_line pil ON p.id = pil.product_id AND pil.demo_item = copl.demo
         JOIN inventory_item ii ON pil.inventory_item_id = ii.id
         JOIN dropbox d ON ii.dropbox_id = d.object_id
		WHERE customer_order.id = $1`,
		orderID)
	if err != nil {
		return nil, err
	}
	var orderItems []OrderItem
	for rows.Next() {
		var orderItem OrderItem
		err = rows.Scan(&orderItem.ObjectID, &orderItem.Name, &orderItem.DropboxURL, &orderItem.DatabaseID, &orderItem.Copyable)
		if err != nil {
			return nil, err
		}
		orderItems = append(orderItems, orderItem)
	}
	return orderItems, nil
}

// ProductOrderItems gets the items pertaining to an order and a specific product.
func ProductOrderItems(ctx context.Context, conn *pgxpool.Pool, orderID string, productID string) ([]OrderItem, error) {
	rows, err := conn.Query(ctx, `
		SELECT ii.object_id, ii.name, d.url, ii.copyable
		FROM customer_order
         JOIN customer_order_product_line copl ON customer_order.id = copl.customer_order
         JOIN product p ON copl.product_id = p.id
         JOIN product_inventory_line pil ON p.id = pil.product_id AND pil.demo_item = copl.demo
         JOIN inventory_item ii ON pil.inventory_item_id = ii.id
         JOIN dropbox d ON ii.dropbox_id = d.object_id
		WHERE customer_order.id = $1 AND p.id = $2`,
		orderID, productID)
	if err != nil {
		return nil, err
	}
	var orderItems []OrderItem
	for rows.Next() {
		var orderItem OrderItem
		err = rows.Scan(&orderItem.ObjectID, &orderItem.Name, &orderItem.DropboxURL, &orderItem.Copyable)
		if err != nil {
			return nil, err
		}
		orderItems = append(orderItems, orderItem)
	}
	return orderItems, nil
}
