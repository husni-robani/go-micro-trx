package rpc_server

import (
	"errors"
	"log"
	"net/rpc"
	"task-service/cmd/config"
	"task-service/data"
)

type RPCServer struct {
	App *config.Config
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

// Handlers
func (r RPCServer) CreateTask(payload CreateTaskPayload, result *RPCResponsePayload) error{
	newTask := data.Task{
		Type: "transaction",
		Data: data.Transaction{
			Amount: int64(payload.Amount),
			DebitAccount: payload.DebitAccount,
			CreditAccount: payload.CreditAccount,
		},
	}

	if err := r.App.Models.Task.CreateTask(newTask); err != nil {
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

func (r RPCServer) RejectTask(payload RejectTaskPayload, result *RPCResponsePayload) error {
	task, err := r.App.Models.Task.GetTaskByID(payload.ID)
	if err != nil {
		log.Println("Failed get task: ", err)
		*result = RPCResponsePayload{
			Error: true,
			Message: err.Error(),
		}
		return err
	}

	if task.Status == 1 {
		log.Printf("Failed to reject transaction, task has been approved")
		*result = RPCResponsePayload{
			Error: true,
			Message: errors.New("task has been rejected").Error(),
		}
		return errors.New("task has been rejected")
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

func (r RPCServer) ApproveTask(payload ApproveTaskPayload, result *RPCResponsePayload) error {
	task, err := r.App.Models.Task.GetTaskByID(payload.ID)
	if err != nil {
		log.Println("Failed to get task: ", err)
		*result = RPCResponsePayload{
			Error: true,
			Message: err.Error(),
		}
		return err
	}

	if task.Status != 0 || task.Step != 2 {
		log.Printf("Failed to approve task %v, task has been rejected", task.TaskID)
		*result = RPCResponsePayload{
			Error: true,
			Message: errors.New("task has been rejected").Error(),
		}
		return errors.New("task has been rejected")
	}

	// update task to approve
	if err := task.ApproveTask(); err != nil {
		log.Println("Failed to approve task: ", err)
		*result = RPCResponsePayload{
			Error: true,
			Message: err.Error(),
		}
		return err
	}

	// make and start transaction
	if err := r.makeAndStartTransaction(*task); err != nil {
		log.Println("Failed to start transaction: ", err)
		*result = RPCResponsePayload{
			Error: true,
			Message: err.Error(),
		}
		return err
	}
	

	*result = RPCResponsePayload{
		Error: false,
		Message: "Task approved!",
	}
	
	return nil
}

func (r RPCServer) makeAndStartTransaction(task data.Task) error {
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