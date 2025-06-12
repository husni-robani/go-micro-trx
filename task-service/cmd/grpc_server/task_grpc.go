package grpc_server

import (
	"context"
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
		ts.createTransactionTask(req)

	default:
		return ts.writeErrorResponse(status.Errorf(codes.InvalidArgument, "type '%s' is not available", req.Type))
	}
	
	return ts.writeResponse(false, "task created!")
}

func (ts TaskGRPCServer) createTransactionTask(task *taskpb.CreateTaskRequest) error {
	dataTask := task.Data
	dataTaskContent := task.Data.Content

	log.Printf("task: %v\ndataTask: %v\ndataTaskContent: %v\n", task, dataTask, dataTaskContent)

	return nil
}