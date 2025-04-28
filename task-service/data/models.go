package data

import (
	"database/sql"
)

type Models struct {
	Task Task
}

var db *sql.DB

func New(dbConn *sql.DB) Models {
	db = dbConn

	return Models{}
}

type Task struct {
	TaskID int
	Type string
	Data Transaction
	Status int
	Step int
}

type Transaction struct {
	Amount int64 `json:"amount,omitempty"`
	Status int `json:"status,omitempty"`
	TaskID int `json:"task_id,omitempty"`
	DebitAccount string `json:"debit_account,omitempty"`
	CreditAccount string `json:"credit_account,omitempty"`
}