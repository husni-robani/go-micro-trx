package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func sendLog(name string, data string) error {
	payload := struct{
		Name string `json:"name"`
		Data string	`json:"data"`
	}{
		Name: name,
		Data: data,
	}

	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to marshaling payload: ", err)
		return err
	}

	url := "http://logger-service/log"

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if  err != nil {
		log.Println("Failed to make request: ", err)
		return err
	}

	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		log.Println("Failed to make request: ", err)
		return err
	}

	if res.StatusCode != http.StatusAccepted {
		log.Printf("send log failed with %d status code", res.StatusCode)
		return errors.New("failed to create log")
	}

	return nil
}