package db

import (
	"context"
	"fmt"
	"time"
)

type Notification struct {
	ID           string    `json:"id"`
	CreationDate time.Time `json:"creation_date"`
	AvatarOwner  string    `json:"avatar_owner"`
	Message      string    `json:"message"`
}

func GetNotifications(ctx context.Context, conn Conn, avatarOwnerID string) ([]Notification, error) {
	rows, err := conn.Query(ctx, "SELECT * FROM notification WHERE owner = $1", avatarOwnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		err = rows.Scan(&n.ID, &n.AvatarOwner, &n.CreationDate, &n.Message)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func CreateNotification(ctx context.Context, conn Conn, avatarOwnerID, message string) error {
	fmt.Println("Creating a notification", avatarOwnerID, message)
	_, err := conn.Exec(ctx, `INSERT INTO notification VALUES(uuid_generate_v4(), $1, now()::timestamp, $2)`,
		avatarOwnerID, message)
	if err != nil {
		return err
	}
	return nil
}
