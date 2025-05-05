package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"transaction-service/data"

	"math/rand"
)

func (app *Config) createTransaction(w http.ResponseWriter, r *http.Request) {
	var newTransaction data.Transaction

	if err := app.readJson(r, &newTransaction); err != nil {
		log.Println("failed to read request body")
		app.errorResponse(w, http.StatusBadRequest, errors.New("invalid body"))
		return
	}

	// run transaction (create)
	fmt.Println("Transaction is on process ...")
	isTransactionSuccess := ExecTransaction(&newTransaction)

	if isTransactionSuccess {
		fmt.Println("Transaction Successful!")
		newTransaction.Status = 1
	}else {
		fmt.Println("Transaction Failed!")
		newTransaction.Status = 0
	}

	// add transaction to database
	if err := app.Models.Transaction.Create(newTransaction); err != nil {
		app.errorResponse(w, http.StatusInternalServerError, errors.New("create transaction failed"))
		return
	}

	// send log when transaction success
	if isTransactionSuccess {
		err := SendLog("transaction", fmt.Sprintf("Transaction success to %s with amount %d", newTransaction.CreditAccount, newTransaction.Amount))
		if err != nil {
			app.errorResponse(w, http.StatusInternalServerError, errors.New("failed to send transaction log"))
			return
		}
	}

	// send response
	var responsePayload jsonResponse
	responsePayload.Error = false
	if isTransactionSuccess {
		responsePayload.Message = "transaction failed"
	} else {
		responsePayload.Message = "transaction transaction success"
	}

	app.writeResponse(w, http.StatusAccepted, responsePayload)
}

func ExecTransaction(transaction *data.Transaction) bool {

	// random status
	transaction.Status = rand.Intn(2)

	return transaction.Status == 1
}

func SendLog(name string, data string) error {
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