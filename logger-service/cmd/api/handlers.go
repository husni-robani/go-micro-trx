package main

import (
	"log"
	"logger-service/data"
	"net/http"
)


type RequestPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) writeLog(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
	err := app.readJson(r, &payload)
	if err != nil {
		log.Println("failed to read request body: ", err)
		app.errorResponse(w, http.StatusBadRequest, err)
		return
	}

	var newLog data.LogEntry
	newLog.Name = payload.Name
	newLog.Data = payload.Data

	if err := newLog.InsertOne(newLog); err != nil {
		app.errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	var responsePayload jsonResponse
	responsePayload.Error = false
	responsePayload.Message = "insert log successful!"
	responsePayload.Data = newLog

	app.writeResponse(w, http.StatusCreated, responsePayload)
}