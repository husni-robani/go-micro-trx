package main

import (
	"errors"
	"net/http"
)

func (app *Config) getAllTask(w http.ResponseWriter, r *http.Request) {
	tasks, err := app.Models.Task.GetAll()
	if err != nil {
		app.errorResponse(w, http.StatusInternalServerError, errors.New("failed to get tasks"))
		return
	}

	var responsePayload jsonResponse
	responsePayload.Error = false
	responsePayload.Message = "Get all task successful"
	responsePayload.Data = tasks

	app.writeResponse(w, http.StatusAccepted, responsePayload)
}