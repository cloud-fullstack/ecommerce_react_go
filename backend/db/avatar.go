package db

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

func MiddlewareDB(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("dbConn", pool)
		c.Next()
	}
}

func InsertAvatar(ctx context.Context, conn *pgxpool.Pool, uuid, legacyName string) error {
	_, err := conn.Exec(ctx, "INSERT INTO avatar VALUES($1, $2)", uuid, legacyName)
	if err != nil {
		return err
	}
	return nil
}

func UpdateHUDURL(ctx context.Context, conn *pgxpool.Pool, url, avatarKey string) error {
	_, err := conn.Exec(ctx, "UPDATE avatar SET hud_url = $1, hud_url_last_updated = to_timestamp($2) WHERE uuid=$3", url, time.Now().Unix(), avatarKey)
	if err != nil {
		return err
	}
	return nil
}

func AvatarLegacyName(ctx context.Context, conn *pgxpool.Pool, avatarKey string) (string, error) {
	var name string
	err := conn.QueryRow(ctx, `SELECT legacyname FROM avatar WHERE uuid = $1`, avatarKey).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}
