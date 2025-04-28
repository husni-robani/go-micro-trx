package data

import (
	"context"
	"database/sql"
	"log"
	"time"
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
	Data []byte
	Status int
	Step int
}


func (t *Task) CreateTask(newTask Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
	defer cancel()

	result, err := db.ExecContext(ctx, "INSERT INTO tasks (type, data, status, step) values ($1, $2, 0, 1)", newTask.Type, newTask.Data)
	if err != nil {
		log.Println("Failed to exec query insert: ", err)
		return err
	}

	rowAffected, _ := result.RowsAffected()
	log.Printf("Data inserted | Rows affected: %d", rowAffected)
	return nil
}