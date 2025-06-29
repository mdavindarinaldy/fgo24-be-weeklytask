package models

import (
	"backend3/utils"
	"context"
	"time"
)

func AddToBlacklist(token string, expiresAt time.Time) error {
	conn, err := utils.DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Exec(context.Background(),
		`INSERT INTO blacklist_tokens (token, expires_at) VALUES ($1, $2)`,
		token, expiresAt)
	if err != nil {
		return err
	}
	return nil
}

func IsTokenBlacklisted(token string) (bool, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return false, err
	}
	defer conn.Close()

	var exists bool
	err = conn.QueryRow(context.Background(),
		`SELECT EXISTS (SELECT 1 FROM blacklist_tokens WHERE token = $1 AND expires_at > $2)`,
		token, time.Now()).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
