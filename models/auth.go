package models

import (
	"backend3/utils"
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
)

func CheckUser(email string) (User, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return User{}, err
	}
	defer conn.Close()
	rows, err := conn.Query(context.Background(),
		`SELECT * FROM users WHERE email=$1`, email)

	if err != nil {
		return User{}, err
	}
	user, err := pgx.CollectOneRow[User](rows, pgx.RowToStructByName)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func HandleRegister(user User) error {
	if user.Email == "" || user.Name == "" || user.Password == "" || user.PhoneNumber == "" || user.Pin == "" {
		return errors.New("user data should not be empty")
	}

	conn, err := utils.DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Exec(
		context.Background(),
		`
		INSERT INTO users (name, email, phone_number, password, pin)
		VALUES
		($1,$2,$3,$4,$5);
		`,
		user.Name, user.Email, user.PhoneNumber, user.Password, user.Pin)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			return errors.New("email already used by another user")
		}
		return err
	}
	return nil
}
