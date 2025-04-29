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
		app.taskApprove(w, request_payload.Task)
	case "task-reject":
		app.taskReject(w, request_payload.Task)
	case "task-all":
		app.taskAll(w)

	default:
		log.Println("invalid handle action")
		app.errorResponse(w, http.StatusBadRequest, errors.New("invalid action"))
		return
	}
}

func (app *Config) taskCreate(w http.ResponseWriter, task Task) {
	task.Data.Status = 0
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
		app.errorResponse(w, http.StatusInternalServerError, errors.New("failed to make request"))
		return
	}
	defer resp.Body.Close()

	// check response code from task service
	if resp.StatusCode != http.StatusAccepted {
		log.Printf("create task failed with %d status code", resp.StatusCode)
		app.errorResponse(w, resp.StatusCode, errors.New("create task failed"))
		return
	}

	// send response
	var responsePayload jsonResponse
	responsePayload.Error = false
	responsePayload.Data = task
	responsePayload.Message = "Create task successful"

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
	requestPayload, err := json.Marshal(task)
	if err != nil {
		log.Println("Failed to marshal task: ", err)
		app.errorResponse(w, http.StatusBadRequest, errors.New("invalid body"))
		return
	}

	// make request
	url := "http://task-service/reject"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestPayload))
	if err != nil {
		log.Println("Failed to make request: ", err)
		app.errorResponse(w, http.StatusInternalServerError, errors.New("failed to reject task"))
		return
	}

	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		log.Println("Failed to make request: ", err)
		app.errorResponse(w, http.StatusInternalServerError, errors.New("failed to reject task"))
		return
	}
	res.Body.Close()

	// check response status code
	if res.StatusCode != http.StatusAccepted {
		log.Printf("reject task failed with status code %d", res.StatusCode)
		app.errorResponse(w, res.StatusCode, errors.New("reject task failed"))
		return
	}

	// send response
	var responsePayload jsonResponse
	responsePayload.Error = false
	responsePayload.Message = "reject task successful"
	app.writeResponse(w, http.StatusOK, responsePayload)
}

func (app *Config) taskAll(w http.ResponseWriter) {
	url := "http::/task-service/all"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("failed to make request: ", err)
		app.errorResponse(w, http.StatusBadRequest, errors.New("invalid request"))
		return
	}
	
	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		log.Println("failed to make request: ", err)
		app.errorResponse(w, http.StatusBadRequest, errors.New("invalid request"))
		return
	}
	defer res.Body.Close()

	// read all tasks from response body
	var responsePayload jsonResponse
	decoder := json.NewDecoder(res.Body)
	
	if err := decoder.Decode(&responsePayload); err != nil {
		log.Println("Failed to decode response body: ", err)
		app.errorResponse(w, res.StatusCode, errors.New("internal server error"))
		return
	}

	// send response
	var payload jsonResponse
	payload.Error = false
	payload.Message = "get tasks successful"
	payload.Data = responsePayload.Data

	app.writeResponse(w, http.StatusOK, payload)
}