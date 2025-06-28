package models

import (
	"backend3/utils"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type TopUpRequest struct {
	Nominal float64 `form:"nominal" json:"nominal" binding:"required"`
}

type TransferRequest struct {
	Nominal     float64 `form:"nominal" json:"nominal" binding:"required"`
	OtherUserId int     `form:"otherUserId" json:"otherUserId" binding:"required"`
	Notes       string  `form:"notes" json:"notes"`
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
	MakeAccountBalance(userId, newBalance)
	return nil
}

func HandleTransfer(request TransferRequest, userId int) error {
	conn, err := utils.DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()
	senderBalance := GetLatestBalance(userId)
	if senderBalance < request.Nominal {
		return errors.New("insufficient balance")
	}
	_, err = conn.Exec(context.Background(),
		`
		INSERT INTO transactions (nominal, type, id_user, id_other_user, notes) 
		VALUES ($1,'expense',$2,$3,$4)
		`, request.Nominal, userId, request.OtherUserId, request.Notes)
	if err != nil {
		return err
	}
	newSenderBalance := senderBalance - request.Nominal
	MakeAccountBalance(userId, newSenderBalance)

	receiverBalance := GetLatestBalance(request.OtherUserId)
	newReceiverBalance := receiverBalance + request.Nominal
	MakeAccountBalance(request.OtherUserId, newReceiverBalance)
	return nil
}
