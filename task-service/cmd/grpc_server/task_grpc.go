package grpc_server

import (
	"context"
	"fmt"
	"log"
	taskpb "proto/task"
	transactionpb "proto/transaction"
	"task-service/data"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskGRPCServer struct {
	taskpb.UnimplementedTaskServiceServer
	Models data.Models
	GRPCClientTransaction transactionpb.TransactionServiceClient
	
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

func (ts TaskGRPCServer) ApproveTask(ctx context.Context, req *taskpb.ApproveTaskRequest) (*taskpb.TaskResponse, error) {
	task, err := ts.Models.Task.GetTaskByID(int(req.TaskID))
	if err != nil {
		return ts.writeResponse(true, "Task Not Found!")
	}

	if task.Status != 0 || task.Step != 2 {
		log.Printf("Failed to approve task %v, task has been rejected", task.TaskID)
		return ts.writeResponse(true, "task has been rejected")
	}

	// update task to approve
	if err := task.ApproveTask(); err != nil {
		log.Println("Failed to approve task: ", err)
		return ts.writeResponse(true, err.Error())
	}

	// make and start transaction
	if err := ts.makeAndStartTransaction(*task); err != nil {
		return ts.writeResponse(true, err.Error())
	}
	
	
	return ts.writeResponse(false, "task approved!")
}


func (ts TaskGRPCServer) makeAndStartTransaction(task data.Task) error {
	requestPayload := transactionpb.CreateTransactionRequest{
		TaskID: int32(task.TaskID),
		DebitAccount: task.Data.DebitAccount,
		CreditAccount: task.Data.CreditAccount,
		Amount: task.Data.Amount,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 3)
	defer cancel()

	res, err := ts.GRPCClientTransaction.CreateTransaction(ctx, &requestPayload)
	if err != nil {
		log.Println("gRPC | Failed to call method CreateTransaction: ", err)
		return err
	}

	log.Printf("CreateTransaction Response: %v", res)
	return nil
}

func (ts TaskGRPCServer) RejectTask(ctx context.Context, req *taskpb.RejectTaskRequest) (*taskpb.TaskResponse, error) {
	task, err := ts.Models.Task.GetTaskByID(int(req.TaskID))
	if err != nil {
		log.Println("Failed get task: ", err)
		return ts.writeResponse(true, fmt.Sprintf("Task not found: %s", err))
	}

	if task.Status == 1 {
		log.Printf("Failed to reject transaction, task has been approved")
		return ts.writeResponse(true, "task has been rejected")
	}

	if err := task.RejectTask(); err != nil {
		log.Println("Failed to reject task: ", err)
		return ts.writeResponse(true, "failed to reject task: " + err.Error())
	}

	return ts.writeResponse(false, "task rejected!")
}