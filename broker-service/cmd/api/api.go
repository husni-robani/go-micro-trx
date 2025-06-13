package api

import (
	"broker-service/cmd/config"
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	taskpb "proto/task"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type APIHandler struct {
	App *config.Config
}

type RequestPayload struct {
	Action string `json:"action"`
	Task Task  `json:"task,omitempty"` 
}

type Task struct {
	TaskID int `json:"task_id,omitempty"`
	Type string `json:"type,omitempty"`
	Data Transaction `json:"data,omitempty"`
	Status int `json:"status,omitempty"`
	Step int `json:"step,omitempty"`
}

type Transaction struct {
	Amount int64 `json:"amount,omitempty"`
	Status int `json:"status,omitempty"`
	DebitAccount string `json:"debit_account,omitempty"`
	CreditAccount string `json:"credit_account,omitempty"`
}

type RPCResponsePayload struct {
	Error bool `json:"error"`
	Message string `json:"message,omitempty"`
}

		// RPC Implementation
type CreateTaskPayload struct {
	Amount int
	DebitAccount string
	CreditAccount string
}

type RejectTaskPayload struct {
	ID int
}

type ApproveTaskPayload struct {
	ID int
}


func (api *APIHandler) Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300,
	}))
	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/handle", api.handleSubmition)

	return mux
}

// handlers
func (api *APIHandler) handleSubmition(w http.ResponseWriter, r *http.Request) {
	var request_payload RequestPayload

	if err := api.readJson(r, &request_payload); err != nil {
		log.Println("Failed to read body request: ", err)
		api.errorResponse(w, http.StatusBadRequest, errors.New("invalid body request"))
		return
	}

	switch request_payload.Action {
	case "rpc-task-create":
		api.taskCreateViaRPC(w, request_payload.Task)
	case "grpc-task-create":
		api.taskCreateViaGRPC(w, request_payload.Task)
	case "rpc-task-approve":
		api.taskApproveViaRPC(w, request_payload.Task)
	case "grpc-task-approve":
		api.taskApproveViaGRPC(w, request_payload.Task)
	case "task-reject":
		api.taskReject(w, request_payload.Task)

	default:
		log.Println("invalid handle action")
		api.errorResponse(w, http.StatusBadRequest, errors.New("invalid action"))
		return
	}
}

func (api *APIHandler) taskCreateViaGRPC(w http.ResponseWriter, task Task) {
	newTaskrequestPayload := taskpb.CreateTaskRequest{
		Type: "transaction",
		Data: &taskpb.DataTask{
			Content: &taskpb.DataTask_Transaction{
				Transaction: &taskpb.TransactionDataTask{
					Amount: task.Data.Amount,
					DebitAccount: task.Data.DebitAccount,
					CreditAccount: task.Data.CreditAccount,
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
	defer cancel()

	grpcRes, err := api.App.GRPCClientTask.CreateTask(ctx, &newTaskrequestPayload)
	if err != nil {
		log.Println("gRPC | CreateTask error: ", err)
		api.errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	var responsePayload jsonResponse
	responsePayload.Error = false
	responsePayload.Message = grpcRes.Message
	api.writeResponse(w, http.StatusCreated, responsePayload)
}

func (api *APIHandler) taskCreateViaRPC(w http.ResponseWriter, task Task) {
	var rpcResponse RPCResponsePayload
	payload := CreateTaskPayload {
		Amount: int(task.Data.Amount),
		DebitAccount: task.Data.DebitAccount,
		CreditAccount: task.Data.CreditAccount,
	}

	rpcMethod := "RPCServer.CreateTask"
	if err := api.App.RPCClientTask.Call(rpcMethod, &payload, &rpcResponse); err != nil {
		log.Printf("error while call %v: %v\n", rpcMethod, err)
		api.errorResponse(w, http.StatusInternalServerError, err)
		return
	}


	var responsePayload jsonResponse
	responsePayload.Error = false
	responsePayload.Message = rpcResponse.Message
	api.writeResponse(w, http.StatusCreated, responsePayload)
}

func (api *APIHandler) taskApproveViaGRPC(w http.ResponseWriter, task Task) {
	grpc_payload := taskpb.ApproveTaskRequest{
		TaskID: int32(task.TaskID),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 3)
	defer cancel()

	res, err := api.App.GRPCClientTask.ApproveTask(ctx, &grpc_payload)
	if err != nil {
		log.Println("gRPC | ApproveTask error: ", err)
		api.errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	var responsePayload jsonResponse
	responsePayload.Error = false
	responsePayload.Message = res.Message
	api.writeResponse(w, http.StatusOK, responsePayload)
}

func (api *APIHandler) taskApproveViaRPC(w http.ResponseWriter, task Task) {
	var rpcResponse RPCResponsePayload
	payload := ApproveTaskPayload{
		task.TaskID,
	}
	
	rpcMethod := "RPCServer.ApproveTask"
	if err := api.App.RPCClientTask.Call(rpcMethod, payload, &rpcResponse); err != nil {
		log.Printf("error while call %v: %v\n", rpcMethod, err)
		api.errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	// send success response
	var responsePayload jsonResponse
	responsePayload.Error = false
	responsePayload.Message = "approve task successful"

	api.writeResponse(w, http.StatusOK, responsePayload)
}


func (api *APIHandler) taskReject(w http.ResponseWriter, task Task) {
	var rpcResponse RPCResponsePayload
	payload := RejectTaskPayload{
		ID: task.TaskID,
	}
	rpcMethod := "RPCServer.RejectTask"

	if err := api.App.RPCClientTask.Call(rpcMethod, &payload, &rpcResponse); err != nil {
		log.Printf("error while call %v: %v\n", rpcMethod, err)
		api.errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	var responsePayload jsonResponse
	responsePayload.Error = false
	responsePayload.Message = "task Rejected!"
	api.writeResponse(w, http.StatusOK, responsePayload)
}