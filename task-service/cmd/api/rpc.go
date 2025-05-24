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

type CreateTaskPayload struct {
	Amount int
	DebitAccount string
	CreditAccount string
}

type RejectTaskPayload struct {
	ID int
}

func NewRPCServer(models data.Models) TaskRPCServer {
	return TaskRPCServer{
		Models: models,
	}
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


// Make implementation of approve and reject task in RPC

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