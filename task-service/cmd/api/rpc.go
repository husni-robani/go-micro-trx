package main

import (
	"log"
	"task-service/data"
)

type RPCServer struct {
	Models data.Models
}

type RPCResponsePayload struct {
	Error bool `json:"error"`
	Data any `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewRPCServer(models data.Models) RPCServer {
	return RPCServer{
		Models: models,
	}
}

type CreatePayload struct {
	Amount int
	DebitAccount string
	CreditAccount string
}

func (r *RPCServer) CreateTask(payload CreatePayload, result *RPCResponsePayload) error{
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
			Data: nil,
			Message: err.Error(),
		}
	}

	*result = RPCResponsePayload{
		Error: false,
		Data: newTask,
		Message: "Make task successful",
	}
	return nil
}