package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"transaction-service/cmd/publisher"
	"transaction-service/data"
)

type TransactionRPCServer struct {
	models data.Models
	publisher publisher.Publisher
}

func NewTransactionRPCServer(m data.Models, p publisher.Publisher) TransactionRPCServer {
	return TransactionRPCServer{
		models: m,
		publisher: p,
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
	

	if isTransactionSuccess {
		var wg sync.WaitGroup
		var mu sync.Mutex
		var errorPublish []error

		wg.Add(2)
		go func(){
			defer wg.Done()
			
			// log message
			logMessage := publisher.LogMessage{
				Name: "transaction",
				Data: fmt.Sprintf("Transaction success to %s with amount %d", newTransaction.CreditAccount, newTransaction.Amount),
			}
									// err := rpc.LogPublisher.PublishMessage(jsonLogMessage)
			err := rpc.publisher.PublishLogMessage(logMessage)
			if err != nil {
				mu.Lock()
				errorPublish = append(errorPublish, err)
				mu.Unlock()
			}
		}()

		go func(){
			defer wg.Done()

			// email message
			mailMessage := publisher.MailMessage{
				Recipient: "husnir2005@gmail.com",
				DebitAccount: newTransaction.DebitAccount,
				CreditAccount: newTransaction.CreditAccount,
				Amount: int64(newTransaction.Amount),
			}
									// err := rpc.MailPublisher.PublishMessage(jsonMailMessage)
			err := rpc.publisher.PublishMailMessage(mailMessage)
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