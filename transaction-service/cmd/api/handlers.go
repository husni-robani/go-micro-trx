package main

import (
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