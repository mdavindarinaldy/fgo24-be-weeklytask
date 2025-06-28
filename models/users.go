package models

import (
	"backend3/utils"
	"context"
	"errors"
	"strings"
)

type User struct {
	Id          int    `db:"id"`
	Name        string `form:"name" json:"name" db:"name" binding:"required"`
	Email       string `form:"email" json:"email" db:"email" binding:"required,email"`
	PhoneNumber string `form:"phoneNumber" json:"phoneNumber" db:"phone_number" binding:"required"`
	Password    string `form:"password" json:"password" db:"password" binding:"required"`
	Pin         string `form:"pin" json:"pin" db:"pin" binding:"required"`
}

func HandleUpdate(user User, id int) error {
	if user.Email == "" || user.Name == "" || user.Password == "" || user.PhoneNumber == "" || user.Pin == "" {
		return errors.New("user data should not be empty")
	}

	conn, err := utils.DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Exec(context.Background(), `
			UPDATE users SET 
			email = $1, 
			name = $2, 
			password = $3, 
			phone_number = $4, 
			pin = $5 
			WHERE id = $6
		`, user.Email, user.Name, user.Password, user.PhoneNumber, user.Pin, id)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			return errors.New("email already used by another user")
		}
		return err
	}
	return nil
}

func GetUsers() {}
