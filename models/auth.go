package models

import (
	"backend3/utils"
	"context"
	"errors"
	"strings"
)

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
	check, _ := GetUserByEmail(user.Email)
	MakeAccountBalance(check.Id, float64(0))
	return nil
}
