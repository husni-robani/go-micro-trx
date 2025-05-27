package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type jsonResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
}

func (app *APIHandler) writeResponse(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed to marshaling data: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Write(jsonData)
	return nil
}


func (app *APIHandler) errorResponse(w http.ResponseWriter, statusCode int, err error) {
	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()
	
	app.writeResponse(w, statusCode, payload)
}