package data

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Transaction struct {
	TransactionID int `json:"transaction_id,omitempty"`
	TaskID int `json:"task_id"`
	DebitAccount string `json:"debit_account"`
	CreditAccount string `json:"credit_account"`
	Amount int `json:"amount"`
	Status int `json:"status,omitempty"`
}

type Models struct {
	Transaction Transaction
}

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{}
}

func (t Transaction) Create(newTransaction Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
	defer cancel()

	query := "INSERT INTO transactions (task_id, debit_account, credit_account, amount, status) VALUES ($1, $2, $3, $4, $5)"

	result, err := db.ExecContext(ctx, query, newTransaction.TaskID, newTransaction.DebitAccount, newTransaction.CreditAccount, newTransaction.Amount, newTransaction.Status)
	if err != nil {
		log.Println("failed to insert new transaction: ", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()

	log.Println("Insert new transaction successful | rows affected: ", rowsAffected)
	return nil
}