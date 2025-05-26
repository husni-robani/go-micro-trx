package main

import (
	"errors"
	"log"
	"net/rpc"
	"task-service/data"
)

type TaskRPCServer struct {
	Models data.Models
	RPCClientTransaction *rpc.Client
}

type RPCResponsePayload struct {
	Error bool `json:"error"`
	Message string `json:"message,omitempty"`
}

type CreateTaskPayload struct {
	Amount int
	DebitAccount string
	CreditAccount string
}

type CreateTransactionPayload struct {
	TaskID int `json:"task_id"`
	DebitAccount string `json:"debit_account"`
	CreditAccount string `json:"credit_account"`
	Amount int64 `json:"amount"`
}

type RejectTaskPayload struct {
	ID int
}

type ApproveTaskPayload struct {
	ID int
}

func (r *TaskRPCServer) CreateTask(payload CreateTaskPayload, result *RPCResponsePayload) error{
	newTask := data.Task{
		Type: "transaction",
		Data: data.Transaction{
			Amount: int64(payload.Amount),
			DebitAccount: payload.DebitAccount,
			CreditAccount: payload.CreditAccount,
		},
	}

	if err := r.Models.Task.CreateTask(newTask); err != nil {
		log.Println("Failed to make task: ", err)
		*result = RPCResponsePayload{
			Error: true,
			Message: err.Error(),
		}
		return err
	}

	*result = RPCResponsePayload{
		Error: false,
		Message: "task created!",
	}
	return nil
}

func (r *TaskRPCServer) RejectTask(payload RejectTaskPayload, result *RPCResponsePayload) error {
	task, err := r.Models.Task.GetTaskByID(payload.ID)
	if err != nil {
		log.Println("Failed get task: ", err)
		*result = RPCResponsePayload{
			Error: true,
			Message: err.Error(),
		}
		return err
	}

	if err := task.RejectTask(); err != nil {
		log.Println("Failed to reject task: ", err)
		*result = RPCResponsePayload{
			Error: true,
			Message: err.Error(),
		}
		return err
	}

	*result = RPCResponsePayload{
		Error: false,
		Message: "Task rejected!",
	}
	return nil
}

func (r *TaskRPCServer) ApproveTask(payload ApproveTaskPayload, result *RPCResponsePayload) error {
	task, err := r.Models.Task.GetTaskByID(payload.ID)
	if err != nil {
		log.Println("Failed to get task: ", err)
		*result = RPCResponsePayload{
			Error: true,
			Message: err.Error(),
		}
		return err
	}

	if task.Status != 0 || task.Step != 2 {
		log.Println("Task cannot be approve")
		return errors.New("task cannot be approve")
	}

	// update task to approve
	if err := task.ApproveTask(); err != nil {
		log.Println("Failed to approve task: ", err)
		return err
	}

	// make and start transaction
	if err := r.makeAndStartTransaction(*task); err != nil {
		return err
	}
	

	*result = RPCResponsePayload{
		Error: false,
		Message: "Task approved!",
	}
	
	return nil
}

func (r *TaskRPCServer) makeAndStartTransaction(task data.Task) error {
	method := "TransactionRPCServer.CreateTransaction"
	
	var responsePayload RPCResponsePayload

	trxPayload := CreateTransactionPayload{
		TaskID: task.TaskID,
		DebitAccount: task.Data.DebitAccount,
		CreditAccount: task.Data.CreditAccount,
		Amount: task.Data.Amount,
	}

	if err := r.RPCClientTransaction.Call(method, trxPayload, &responsePayload); err != nil {
		log.Println("Failed to make and start transaction: ", responsePayload.Message)
		return err
	}

	return nil
}