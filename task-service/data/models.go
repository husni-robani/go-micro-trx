package data

import (
	"context"
	"database/sql"
	"encoding/json"
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
	TaskID int `json:"task_id"`
	Type string `json:"type"`
	Data  Transaction `json:"data"`
	Status int `json:"status"`
	Step int `json:"step"`
}

type Transaction struct {
	Amount int64 `json:"amount,omitempty"`
	Status int `json:"status,omitempty"`
	DebitAccount string `json:"debit_account,omitempty"`
	CreditAccount string `json:"credit_account,omitempty"`
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

func (t *Task) ApproveTask() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
	defer cancel()

	query := "UPDATE tasks set status = 1, step = 2 WHERE task_id = $1"

	result, err := db.ExecContext(ctx, query, t.TaskID)
	if err != nil {
		log.Println("Failed to update task to approve: ", err)
		return err
	}

	rowAffected, _ := result.RowsAffected()
	log.Println("Task updated | rows affected: ", rowAffected)

	return nil
}

func (t *Task) RejectTask() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 20)
	defer cancel()

	query := "UPDATE tasks SET status = 2 WHERE task_id = $1"
	result, err := db.ExecContext(ctx, query, t.TaskID)
	if err != nil {
		log.Println("Failed to reject task: ", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Println("task updated to reject | rows affected: ", rowsAffected)

	return nil
}

func (t *Task) GetTaskByID(taskId int) (*Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
	defer cancel()

	var taskResult Task

	row := db.QueryRowContext(ctx, "SELECT task_id, type, data, status, step FROM tasks WHERE task_id = $1", taskId)
	if err :=  row.Scan(&taskResult.TaskID, &taskResult.Type, &taskResult.Data, &taskResult.Status, &taskResult.Step); err != nil {
		log.Println("failed to scan task from database: ", err)
		return nil, err
	}

	return &taskResult, nil
}

func (t *Task) GetAll() ([]Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, "SELECT task_id, type, data, status, step FROM tasks")
	if err != nil {
		log.Println("Failed to exec select query: ", err)
		return nil, err
	}

	var tasks []Task

	for rows.Next(){
		var task Task
		var tempDataTask []uint8

		if err := rows.Scan(&task.TaskID, &task.Type, &tempDataTask, &task.Status, &task.Step); err != nil {
			log.Println("failed to scan task from rows: ", err)
			return nil, err
		}

		// TODO: PERBAIK INI
		var dataTask Transaction
		if err := json.Unmarshal(tempDataTask, &dataTask); err != nil {
			log.Println("failed to unmarshal data: ", err)
			return nil, err
		}

		task.Data = dataTask
		
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}