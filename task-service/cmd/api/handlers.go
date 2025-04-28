package main

import (
	"encoding/json"
	"errors"
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