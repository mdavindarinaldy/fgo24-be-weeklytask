package models

import (
	"backend3/utils"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type TopUpRequest struct {
	Nominal float64 `form:"nominal" json:"nominal" binding:"required"`
}

func MakeAccountBalance(id int, balance float64) error {
	conn, err := utils.DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Exec(context.Background(),
		`
		INSERT INTO user_balance (id_user, balance) VALUES ($1,$2)
		`, id, balance)
	if err != nil {
		return err
	}
	return nil
}

func GetLatestBalance(id int) float64 {
	type userBalance struct {
		Id        int       `db:"id"`
		IdUser    int       `db:"id_user"`
		Balance   float64   `db:"balance"`
		CreatedAt time.Time `db:"created_at"`
	}
	conn, _ := utils.DBConnect()
	defer conn.Close()
	rows, _ := conn.Query(context.Background(),
		`
		SELECT * FROM user_balance WHERE id_user=$1 
		ORDER BY created_at DESC
		LIMIT 1`, id)
	balance, _ := pgx.CollectOneRow[userBalance](rows, pgx.RowToStructByName)
	fmt.Println(balance)
	return balance.Balance
}

func HandleTopUp(request TopUpRequest, userId int) error {
	conn, err := utils.DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Exec(context.Background(),
		`
		INSERT INTO transactions (nominal, type, id_user, id_other_user, notes) 
		VALUES ($1,'income',$2,$3,'top up')
		`, request.Nominal, userId, userId)
	if err != nil {
		return err
	}
	currentBalance := GetLatestBalance(userId)
	newBalance := currentBalance + request.Nominal
	fmt.Println(currentBalance)
	MakeAccountBalance(userId, newBalance)
	return nil
}
