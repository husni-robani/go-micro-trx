package main

import (
	"log"
	"task-service/data"
)

type TaskRPCServer struct {
	Models data.Models
}

type RPCResponsePayload struct {
	Error bool `json:"error"`
	Message string `json:"message,omitempty"`
}

func NewRPCServer(models data.Models) TaskRPCServer {
	return TaskRPCServer{
		Models: models,
	}
}

type CreatePayload struct {
	Amount int
	DebitAccount string
	CreditAccount string
}

func (r *TaskRPCServer) CreateTask(payload CreatePayload, result *RPCResponsePayload) error{
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