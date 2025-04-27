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
	TaskID int `json:"task_id,omitempty"`
	DebitAccount string `json:"debit_account,omitempty"`
	CreditAccount string `json:"credit_account,omitempty"`
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
	case "task-reject":
	default:
		log.Println("invalid handle action")
		app.errorResponse(w, http.StatusBadRequest, errors.New("invalid action"))
		return
	}
}

func (app *Config) taskCreate(w http.ResponseWriter, task Task) {
	requestPayload, err := json.Marshal(task)
	if err != nil {
		log.Println("Failed to marshal task: ", err)
		app.errorResponse(w, http.StatusBadRequest, errors.New("invalid body"))
		return
	}

	// make a post request to /create
	url := "http://task-service/create"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestPayload))
	if err != nil {
		log.Println("Failed to make request: ", err)
		app.errorResponse(w, http.StatusInternalServerError, errors.New("failed to make request"))
		return
	}

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("Failed to make request: ", err)
		app.errorResponse(w, http.StatusBadRequest, errors.New("invalid to make request"))
		return
	}
	defer resp.Body.Close()

	// check response code from task service
	if resp.StatusCode != http.StatusAccepted {
		log.Printf("create task failed with %d status code: %v", resp.StatusCode, err)
		app.errorResponse(w, http.StatusBadRequest, errors.New("create task failed"))
		return
	}

	// send response
	var responsePayload jsonResponse
	responsePayload.Error = false
	responsePayload.Data = task
	responsePayload.Message = "Create task successful"

	app.writeResponse(w, http.StatusCreated, responsePayload)
}