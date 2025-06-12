package grpc_server

import (
	"context"
	"fmt"
	"log"
	taskpb "proto/task"
	"task-service/data"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskGRPCServer struct {
	taskpb.UnimplementedTaskServiceServer
	Models data.Models
}

func (ts TaskGRPCServer) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.TaskResponse, error) {
	
	switch {
	case req.Type == "transaction":
		err := ts.createTransactionTask(req)
		if err != nil {
			return ts.writeResponse(true, fmt.Sprintf("Failed to add task to database: %v", err))
		}
		return ts.writeResponse(false, "task created!")
	default:
		return ts.writeErrorResponse(status.Errorf(codes.InvalidArgument, "type '%s' is not available", req.Type))
	}

}

func (ts TaskGRPCServer) createTransactionTask(task *taskpb.CreateTaskRequest) error {
	newTask := data.Task{
		Type: task.Type,
		Data: data.Transaction{
			Amount: task.Data.GetTransaction().Amount,
			DebitAccount: task.Data.GetTransaction().DebitAccount,
			CreditAccount: task.Data.GetTransaction().CreditAccount,
		},
	}

	if err := ts.Models.Task.CreateTask(newTask); err != nil {
		log.Println("Failed to add new task to database: ", err)
		return err
	}

	fmt.Printf("POSTGRES: Task Created!\nData: %#v", newTask)

	return nil
}