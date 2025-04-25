package main

import (
	"errors"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string `json:"action"`
	Task CreateTaskPayload `json:"task,omitempty"` 
}

type CreateTaskPayload struct {
	
}

type ApproveTaskPayload struct {

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
	default:
		log.Println("invalid handle action")
		app.errorResponse(w, http.StatusBadRequest, errors.New("invalid action"))
		return
	}
}

func (app *Config) taskCreate(w http.ResponseWriter, task CreateTaskPayload) {
	// create task
}