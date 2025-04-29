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

	result, err := db.ExecContext(ctx, "INSERT INTO tasks (type, data, status, step) values ($1, $2, 0, 2)", newTask.Type, newTask.Data)
	if err != nil {
		log.Println("Failed to exec query insert: ", err)
		return err
	}

	rowAffected, _ := result.RowsAffected()
	log.Printf("Data inserted | Rows affected: %d", rowAffected)
	return nil
}

func (t *Task) ApproveTask(taskId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
	defer cancel()

	query := "UPDATE tasks set status = 1, step = 2 WHERE task_id = $1"

	result, err := db.ExecContext(ctx, query, taskId)
	if err != nil {
		log.Println("Failed to update task to approve: ", err)
		return err
	}

	rowAffected, _ := result.RowsAffected()
	log.Println("Task updated | rows affected: ", rowAffected)

	return nil
}

func (t *Task) RejectTask(taskId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 20)
	defer cancel()

	query := "UPDATE tasks SET status = 2 WHERE task_id = $1"
	result, err := db.ExecContext(ctx, query, taskId)
	if err != nil {
		log.Println("Failed to reject task: ", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Println("task updated to reject | rows affected: ", rowsAffected)

	return nil
}