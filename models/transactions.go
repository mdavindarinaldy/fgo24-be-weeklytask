package models

import (
	"backend3/utils"
	"context"
	"errors"
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

type Transactions struct {
	TransactionsDate     time.Time `db:"transactions_date" json:"transactionDate"`
	Nominal              float64   `db:"nominal" json:"nominal"`
	Type                 string    `db:"type" json:"type"`
	Notes                string    `db:"notes" json:"notes"`
	IdOtherUser          int       `db:"id_other_user" json:"idOtherUser"`
	OtherUserName        string    `db:"other_user_name" json:"otherUserName"`
	OtherUserEmail       string    `db:"other_user_email" json:"otherUserEmail"`
	OtherUserPhoneNumber string    `db:"other_user_phone" json:"otherUserPhone"`
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

func GetHistoryExpenseTransactions(id int, page int) ([]Transactions, PageData, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return []Transactions{}, PageData{}, err
	}
	defer conn.Close()

	type Count struct {
		Count int `db:"count"`
	}
	count, err := conn.Query(context.Background(),
		`SELECT COUNT(*) as count FROM transactions 
		WHERE type='expense' AND id_user=$1`, id)
	if err != nil {
		return []Transactions{}, PageData{}, err
	}
	countData, err := pgx.CollectOneRow[Count](count, pgx.RowToStructByName)
	if err != nil {
		return []Transactions{}, PageData{}, err
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
		`SELECT t.transactions_date, t.nominal, t.type,  
		t.notes, t.id_other_user, 
		u.name AS other_user_name, 
		u.email AS other_user_email, 
		u.phone_number AS other_user_phone 
		FROM transactions t 
		JOIN users u ON u.id = t.id_other_user
		WHERE t.id_user=$1
		ORDER BY t.transactions_date DESC
		OFFSET $2
		LIMIT 5`, id, offset)

	if err != nil {
		return []Transactions{}, PageData{}, err
	}
	transactions, err := pgx.CollectRows[Transactions](rows, pgx.RowToStructByName)
	if err != nil {
		return []Transactions{}, PageData{}, err
	}

	return transactions, pageData, nil
}

func GetHistoryIncomeTransactions(id int, page int) ([]Transactions, PageData, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return []Transactions{}, PageData{}, err
	}
	defer conn.Close()

	type Count struct {
		Count int `db:"count"`
	}
	count, err := conn.Query(context.Background(),
		`SELECT COUNT(*) as count FROM transactions 
		WHERE id_other_user=$1`, id)
	if err != nil {
		return []Transactions{}, PageData{}, err
	}
	countData, err := pgx.CollectOneRow[Count](count, pgx.RowToStructByName)
	if err != nil {
		return []Transactions{}, PageData{}, err
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
		`SELECT 
			t.transactions_date, 
			t.nominal,
		  	CASE 
				WHEN t.type='income' THEN 'income'
				WHEN t.type='expense' AND t.id_other_user=$1 THEN 'income'
			END AS type,
			t.notes, 
			CASE 
				WHEN t.type='income' THEN t.id_user
				WHEN t.type='expense' AND t.id_other_user=$1 THEN t.id_user
			END AS id_other_user, 
			u.name AS other_user_name, 
			u.email AS other_user_email, 
			u.phone_number AS other_user_phone 
		FROM transactions t
		JOIN users u ON 
			(CASE 
				WHEN t.type='income' THEN u.id = t.id_user
				WHEN t.type='expense' AND t.id_other_user=$1 THEN u.id = t.id_user
			END)
		WHERE id_other_user=$1
		ORDER BY transactions_date DESC
		OFFSET $2
		LIMIT 5`, id, offset)

	if err != nil {
		return []Transactions{}, PageData{}, err
	}
	transactions, err := pgx.CollectRows[Transactions](rows, pgx.RowToStructByName)
	if err != nil {
		return []Transactions{}, PageData{}, err
	}

	return transactions, pageData, nil
}

func GetHistoryTransactions(id int, page int) ([]Transactions, PageData, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return []Transactions{}, PageData{}, err
	}
	defer conn.Close()

	type Count struct {
		Count int `db:"count"`
	}
	count, err := conn.Query(context.Background(),
		`SELECT COUNT(*) as count FROM transactions 
		WHERE id_user=$1 OR id_other_user=$1`, id)
	if err != nil {
		return []Transactions{}, PageData{}, err
	}
	countData, err := pgx.CollectOneRow[Count](count, pgx.RowToStructByName)
	if err != nil {
		return []Transactions{}, PageData{}, err
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
		`SELECT 
			t.transactions_date, 
			t.nominal,
			CASE 
				WHEN t.type='income' THEN 'income'
				WHEN t.type='expense' AND t.id_user=$1 THEN 'expense'
				WHEN t.type='expense' AND t.id_other_user=$1 THEN 'income'
			END AS type,
			t.notes,
			CASE 
				WHEN t.type='income' THEN t.id_user
				WHEN t.type='expense' AND t.id_user=$1 THEN t.id_other_user
				WHEN t.type='expense' AND t.id_other_user=$1 THEN t.id_user
			END AS id_other_user,
			u.name AS other_user_name,
			u.email AS other_user_email,
			u.phone_number AS other_user_phone
		FROM transactions t
		JOIN users u ON 
			(CASE 
				WHEN t.type='income' THEN u.id = t.id_user
				WHEN t.type='expense' AND t.id_user=$1 THEN u.id = t.id_other_user
				WHEN t.type='expense' AND t.id_other_user=$1 THEN u.id = t.id_user
			END)
		WHERE t.id_user=$1 OR t.id_other_user=$1
		ORDER BY t.transactions_date DESC
		OFFSET $2
		LIMIT 5`, id, offset)

	if err != nil {
		return []Transactions{}, PageData{}, err
	}
	transactions, err := pgx.CollectRows[Transactions](rows, pgx.RowToStructByName)
	if err != nil {
		return []Transactions{}, PageData{}, err
	}

	return transactions, pageData, nil
}

func GetTotalIncome(id int) (float64, time.Time, time.Time) {
	type Income struct {
		Income float64 `db:"total_income" json:"income"`
	}
	conn, _ := utils.DBConnect()
	defer conn.Close()

	now := time.Now()
	duration := time.Now().Add(-7 * 24 * time.Hour)

	rows, _ := conn.Query(context.Background(),
		`
		SELECT SUM(nominal) AS total_income FROM transactions
		WHERE transactions_date BETWEEN $1 AND $2 
		AND ((type='income' AND id_user=$3) 
		OR (type='expense' AND id_other_user=$3))
		`, duration, now, id)
	income, _ := pgx.CollectOneRow[Income](rows, pgx.RowToStructByName)
	return income.Income, now, duration
}

func GetTotalExpense(id int) (float64, time.Time, time.Time) {
	type Expense struct {
		Expense float64 `db:"total_expense" json:"expense"`
	}
	conn, _ := utils.DBConnect()
	defer conn.Close()

	now := time.Now()
	duration := time.Now().Add(-7 * 24 * time.Hour)

	rows, _ := conn.Query(context.Background(),
		`
		SELECT SUM(nominal) AS total_expense FROM transactions
		WHERE transactions_date BETWEEN $1 AND $2 
		AND type='expense' AND id_user=$3
		`, duration, now, id)
	expense, _ := pgx.CollectOneRow[Expense](rows, pgx.RowToStructByName)
	return expense.Expense, now, duration
}
