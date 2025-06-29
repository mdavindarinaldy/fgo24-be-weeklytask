package models

import (
	"backend3/utils"
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
)

type User struct {
	Id          int    `db:"id" json:"id"`
	Name        string `form:"name" json:"name" db:"name" binding:"required"`
	Email       string `form:"email" json:"email" db:"email" binding:"required,email"`
	PhoneNumber string `form:"phoneNumber" json:"phoneNumber" db:"phone_number" binding:"required"`
	Password    string `form:"password" json:"password" db:"password" binding:"required"`
	Pin         string `form:"pin" json:"pin" db:"pin" binding:"required"`
}

type PageData struct {
	TotalData   int `json:"totalData"`
	TotalPage   int `json:"totalPage"`
	CurrentPage int `json:"currentPage"`
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

func GetUserByEmail(email string) (User, error) {
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

func GetDetailUser(id int) (User, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return User{}, err
	}
	defer conn.Close()
	rows, err := conn.Query(context.Background(),
		`SELECT * FROM users WHERE id=$1`, id)

	if err != nil {
		return User{}, err
	}
	user, err := pgx.CollectOneRow[User](rows, pgx.RowToStructByName)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func GetAllUsers(search string, page int) ([]User, PageData, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return []User{}, PageData{}, err
	}
	defer conn.Close()

	type Count struct {
		Count int `db:"count"`
	}
	count, err := conn.Query(context.Background(),
		`SELECT COUNT(*) as count FROM users 
		WHERE name ILIKE $1 
		OR phone_number ILIKE $1`, "%"+search+"%")
	if err != nil {
		return []User{}, PageData{}, err
	}
	countData, err := pgx.CollectOneRow[Count](count, pgx.RowToStructByName)
	if err != nil {
		return []User{}, PageData{}, err
	}

	offset := (page - 1) * 5
	if page == 0 {
		page = 1
	} else if ((page * 5) - countData.Count) < 5 {
		page = 1
	}

	totalPage := 0
	if countData.Count%5 != 0 {
		totalPage = (countData.Count / 5) + 1
	} else {
		totalPage = countData.Count / 5
	}

	pageData := PageData{
		TotalData:   countData.Count,
		TotalPage:   totalPage,
		CurrentPage: page,
	}

	rows, err := conn.Query(context.Background(),
		`SELECT * FROM users 
		WHERE name ILIKE $1 
		OR phone_number ILIKE $1
		OFFSET $2
		LIMIT 5`, "%"+search+"%", offset)

	if err != nil {
		return []User{}, PageData{}, err
	}
	users, err := pgx.CollectRows[User](rows, pgx.RowToStructByName)
	if err != nil {
		return []User{}, PageData{}, err
	}

	return users, pageData, nil
}
