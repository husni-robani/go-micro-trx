package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
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
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailMessage struct {
	Recipient string `json:"recipient"`
	DebitAccount string `json:"debit_account"`
	CreditAccount string `json:"credit_account"`
	Amount int64 `json:"amount"`
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

	// send log and email when transaction success
		// log message
	jsonLogMessage, err := json.Marshal(LogMessage{
		Name: "transaction",
		Data: fmt.Sprintf("Transaction success to %s with amount %d", newTransaction.CreditAccount, newTransaction.Amount),
	})
	if err != nil {
		log.Println("Failed to marshal json: ", err)
		return err
	}

		// email message
	jsonMailMessage, err := json.Marshal(MailMessage{
		Recipient: "husnir2005@gmail.com",
		DebitAccount: newTransaction.DebitAccount,
		CreditAccount: newTransaction.CreditAccount,
		Amount: int64(newTransaction.Amount),
	})
	if err != nil {
		log.Println("Failed to marshal json: ", err)
		return err
	}

	if isTransactionSuccess {
		var wg sync.WaitGroup
		var mu sync.Mutex
		var errorPublish []error

		wg.Add(2)
		go func(){
			defer wg.Done()

			err := rpc.LogPublisher.PublishMessage(jsonLogMessage)
			if err != nil {
				mu.Lock()
				errorPublish = append(errorPublish, err)
				mu.Unlock()
			}
		}()

		go func(){
			defer wg.Done()

			err := rpc.MailPublisher.PublishMessage(jsonMailMessage)
			if err != nil {
				mu.Lock()
				errorPublish = append(errorPublish, err)
				mu.Unlock()
			}
		}()
		
		wg.Wait()

		if len(errorPublish) > 0 {
			log.Println("error publish message: ", errorPublish)
			*result = RPCResponsePayload{
				Error: true,
				Message: fmt.Sprintf("Publish message failed: %v", errorPublish),
			}
			return errorPublish[0]
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