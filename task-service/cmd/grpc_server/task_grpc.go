package task_grpc_server

import (
	"context"
	"log"
	taskpb "proto/task"
	"task-service/data"
)

type TaskGRPCServer struct {
	taskpb.UnimplementedTaskServiceServer
	Models data.Models
}

func (ts TaskGRPCServer) CreateTask(ctx context.Context, req *taskpb.TaskRequest) (*taskpb.TaskResponse, error) {
	task := req.TaskEntry
	dataTask := req.TaskEntry.Data
	dataTaskContent := req.TaskEntry.Data.Content

	log.Printf("task: %v\ndataTask: %v\ndataTaskContent: %v\n", task, dataTask, dataTaskContent)
	

	var res taskpb.TaskResponse
	res.Result = "task created!"
	return &res, nil
}