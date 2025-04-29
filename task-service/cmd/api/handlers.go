package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"task-service/data"
)

type Transaction struct {
	Amount int64 `json:"amount,omitempty"`
	Status int `json:"status,omitempty"`
	DebitAccount string `json:"debit_account,omitempty"`
	CreditAccount string `json:"credit_account,omitempty"`
}

func (app *Config) createTask(w http.ResponseWriter, r *http.Request) {	
	var requestPayload struct {
		Type string `json:"type"`
		Data Transaction `json:"data"`
	}

	if err := app.readJson(r, &requestPayload); err != nil {
		log.Println("failed to read request payload: ", err)
		app.errorResponse(w, http.StatusBadRequest, err)
		return
	}

	dataMarshal, err := json.Marshal(requestPayload.Data)
	if err != nil {
		log.Println("Failed to marshal data task: ", err)
		app.errorResponse(w, http.StatusBadRequest, err)
		return
	}

	newTask := data.Task{
		Type: requestPayload.Type,
		Data: dataMarshal,
	}

	if err := app.Models.Task.CreateTask(newTask); err != nil {
		app.errorResponse(w, http.StatusInternalServerError, errors.New("failed to create new task"))
		return
	}

	var responsePayload jsonResponse
	responsePayload.Error = false
	responsePayload.Message = "Create new task successful"
	app.writeResponse(w, http.StatusAccepted, responsePayload)
}

func (app *Config) approveTask(w http.ResponseWriter, r *http.Request) {
	// read request
	var requestPayload struct{
		TaskID int `json:"task_id"`
		Data Transaction `json:"data"`
	}

	if err := app.readJson(r, &requestPayload); err != nil {
		log.Println("failed to read request payload: ", err)
		app.errorResponse(w, http.StatusBadRequest, err)
		return
	}

	// Get the task from database to know what step and status of current task
	task, err := app.Models.Task.GetTaskByID(requestPayload.TaskID)
	if err != nil {
		app.errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	if task.Status != 0 || task.Step != 2 {
		app.errorResponse(w, http.StatusBadRequest, errors.New("task cannot be approve"))
		return
	}

	// approve task
	if err := task.ApproveTask(); err != nil {
		app.errorResponse(w, http.StatusInternalServerError, errors.New("failed to approve task"))
		return
	}

	// call transaction-service handler to make transaction
	url := "http://transaction-service/create"
	var sendRequestPayload struct {
		TaskID int `json:"task_id"`
		DebitAccount string `json:"debit_account"`
		CreditAccount string `json:"credit_account"`
		Amount int64 `json:"amount"`
	}
	sendRequestPayload.TaskID = requestPayload.TaskID
	sendRequestPayload.DebitAccount = requestPayload.Data.DebitAccount
	sendRequestPayload.CreditAccount = requestPayload.Data.CreditAccount
	sendRequestPayload.Amount = requestPayload.Data.Amount

	jsonPayload, err := json.Marshal(sendRequestPayload)
	if err != nil {
		log.Println("Failed to marshal send request payload: ", err)
		app.errorResponse(w, http.StatusBadRequest, err)
		return
	}
	
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))

	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		log.Println("Failed to perform post request: ", err)
		app.errorResponse(w, http.StatusInternalServerError, errors.New("failed to make transaction"))
		return
	}
	defer res.Body.Close()

		// check is response code = accepted
	if res.StatusCode != http.StatusAccepted {
		app.errorResponse(w, http.StatusInternalServerError, errors.New("failed to make transaction"))
		return
	}


	// send response
	var responsePayload jsonResponse
	responsePayload.Error = false
	responsePayload.Message = fmt.Sprintf("task %d approved!", requestPayload.TaskID)

	app.writeResponse(w, http.StatusAccepted, responsePayload)
}

func (app *Config) rejectTask(w http.ResponseWriter, r *http.Request){
	// read request body
	var requestPayload struct {
		TaskID int `json:"task_id"`
	}
	if err := app.readJson(r, &requestPayload); err != nil{
		log.Println("Failed to read request body: ", err)
		app.errorResponse(w, http.StatusBadRequest, errors.New("invalid body"))
		return
	}

	// get task from database
	task, err := app.Models.Task.GetTaskByID(requestPayload.TaskID)
	if err != nil {
		app.errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	if task.Status != 0 || task.Step != 2 {
		app.errorResponse(w, http.StatusBadRequest, errors.New("task cannot be reject"))
		return
	}

	// update task
	if err := task.RejectTask(); err != nil {
		app.errorResponse(w, http.StatusInternalServerError, errors.New("failed to reject task"))
		return
	}

	// send response
	var responsePayload jsonResponse
	responsePayload.Error = false
	responsePayload.Message = "task rejected!"
	app.writeResponse(w, http.StatusAccepted, responsePayload)
}