package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

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

type CreateTaskPayload struct {
	Amount int
	DebitAccount string
	CreditAccount string
}

type RejectTaskPayload struct {
	ID int
}

func (app *Config) handleSubmition(w http.ResponseWriter, r *http.Request) {
	var request_payload RequestPayload

	if err := app.readJson(r, &request_payload); err != nil {
		log.Println("Failed to read body request: ", err)
		app.errorResponse(w, http.StatusBadRequest, errors.New("invalid body request"))
		return
	}

	switch request_payload.Action {
	case "task-create":
		app.taskCreate(w, request_payload.Task)
	case "task-approve":
		app.taskApprove(w, request_payload.Task)
	case "task-reject":
		app.taskReject(w, request_payload.Task)

	default:
		log.Println("invalid handle action")
		app.errorResponse(w, http.StatusBadRequest, errors.New("invalid action"))
		return
	}
}

func (app *Config) taskCreate(w http.ResponseWriter, task Task) {
	var rpcResponse RPCResponsePayload
	payload := CreateTaskPayload {
		Amount: int(task.Data.Amount),
		DebitAccount: task.Data.DebitAccount,
		CreditAccount: task.Data.CreditAccount,
	}

	rpcMethod := "TaskRPCServer.CreateTask"
	if err := app.rpcClient.Call(rpcMethod, &payload, &rpcResponse); err != nil {
		log.Printf("Failed to call %v: %v\n", rpcMethod, err)
		app.errorResponse(w, http.StatusInternalServerError, err)
		return
	}


	var responsePayload jsonResponse
	responsePayload.Error = false
	responsePayload.Message = rpcResponse.Message
	app.writeResponse(w, http.StatusCreated, responsePayload)
}

func (app *Config) taskApprove(w http.ResponseWriter, task Task) {
	requestPayload, err := json.Marshal(task)
	if err != nil {
		log.Println("Failed to marshal task: ", err)
		app.errorResponse(w, http.StatusBadRequest, errors.New("invalid body"))
		return
	}
	
	// make a request
	url := "http://task-service/approve"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestPayload))
	if err != nil {
		log.Println("Failed to make request: ", err)
		app.errorResponse(w, http.StatusInternalServerError, errors.New("failed to make request"))
		return
	}

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("Failed to make request: ", err)
		app.errorResponse(w, http.StatusInternalServerError, errors.New("failed to make request"))
		return
	}
	response.Body.Close()

	// check response status code
	if response.StatusCode != http.StatusAccepted {
		log.Printf("approve task failed with %d status code", response.StatusCode)
		app.errorResponse(w, response.StatusCode, errors.New("approve task failed"))
		return
	}

	// send success response
	var responsePayload jsonResponse
	responsePayload.Error = false
	responsePayload.Message = "approve task successful"

	app.writeResponse(w, http.StatusOK, responsePayload)
}


func (app *Config) taskReject(w http.ResponseWriter, task Task) {
	var rpcResponse RPCResponsePayload
	payload := RejectTaskPayload{
		ID: task.TaskID,
	}
	rpcMethod := "TaskRPCServer.RejectTask"

	if err := app.rpcClient.Call(rpcMethod, &payload, &rpcResponse); err != nil {
		log.Println("Failed to call rpc method")
		app.errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	var responsePayload jsonResponse
	responsePayload.Error = false
	responsePayload.Message = "task Rejected!"
	app.writeResponse(w, http.StatusOK, responsePayload)
}