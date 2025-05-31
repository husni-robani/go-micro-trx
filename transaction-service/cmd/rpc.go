package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"transaction-service/data"
)

type TransactionRPCServer struct {
	models data.Models
	LogPublisher Publisher
	MailPublisher Publisher
}

func NewTransactionRPCServer(m data.Models, logPublisher Publisher, mailPublisher Publisher) TransactionRPCServer {
	return TransactionRPCServer{
		models: m,
		LogPublisher: logPublisher,
		MailPublisher: mailPublisher,
	}
}

type RPCResponsePayload struct {
	Error bool 
	Message string
}

type CreateTransactionPayload struct {
	TaskID int
	DebitAccount string
	CreditAccount string
	Amount int64
}

type LogMessage struct {
	Name string
	Data string
}

func (rpc TransactionRPCServer) CreateTransaction(trxPayload CreateTransactionPayload, result *RPCResponsePayload) error {
	newTransaction := data.Transaction{
		TaskID: trxPayload.TaskID,
		DebitAccount: trxPayload.DebitAccount,
		CreditAccount: trxPayload.CreditAccount,
		Amount: int(trxPayload.Amount),
	}
	
	log.Printf("Transaction (task: %d) is on process....", newTransaction.TaskID)

	// exec transaction to get status
	isTransactionSuccess := execTransaction(&newTransaction)

	// add transaction to database
	if err := rpc.models.Transaction.Create(newTransaction); err != nil {
		log.Println("Create new transaction failed: ", err)
		*result = RPCResponsePayload{
			Error: true,
			Message: err.Error(),
		}
		return err
	}

	// send log when transaction success
	jsonMessage, err := json.Marshal(LogMessage{
		Name: "transaction",
		Data: fmt.Sprintf("Transaction success to %s with amount %d", newTransaction.CreditAccount, newTransaction.Amount),
	})
	if err != nil {
		log.Println("Failed to marshal message: ", err)
		return err
	}

	if isTransactionSuccess {
		err := rpc.LogPublisher.PublishMessage(jsonMessage)
		if err != nil {
			log.Println("Failed to send log to mongoDB: ", err)
			*result = RPCResponsePayload{
				Error: true,
				Message: err.Error(),
			}
			return err
		}
	}

	*result = RPCResponsePayload{
		Error: false,
		Message: "Transaction Created!",
	}
	return nil
}

func execTransaction(transaction *data.Transaction) bool {
	// random status
	transaction.Status = rand.Intn(2)

	return transaction.Status == 1
}