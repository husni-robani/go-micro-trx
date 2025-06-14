package grpc_server

import (
	"context"
	"fmt"
	"log"
	transactionpb "proto/transaction"
	"transaction-service/cmd/publisher"
	"transaction-service/cmd/utils"
	"transaction-service/data"
)

type TransactionGRPCServer struct {
	transactionpb.UnimplementedTransactionServiceServer
	models data.Models
	publisher publisher.Publisher
}

func NewTransactionGRPCServer(m data.Models, p publisher.Publisher) TransactionGRPCServer {
	return TransactionGRPCServer{
		models: m,
		publisher: p,
	}
}

func (ts TransactionGRPCServer) CreateTransaction(ctx context.Context, req *transactionpb.CreateTransactionRequest) (*transactionpb.TransactionResponse, error){
	newTransaction := data.Transaction{
		TaskID: int(req.TaskID),
		DebitAccount: req.DebitAccount,
		CreditAccount: req.CreditAccount,
		Amount: int(req.Amount),
	}

	log.Printf("Transaction (TaskID: %v) is on process", newTransaction.TaskID)

	// transaction running process
	if err := utils.ExecuteTransaction(&newTransaction); err != nil {
		return ts.writeResponse(false, err.Error())
	}

	// add to database
	if err := ts.models.Transaction.Create(newTransaction); err != nil {
		log.Println("add transaction to database failed: ", err)
		return ts.writeErrorResponse(err)
	}

	// send message for logging
	go ts.publisher.PublishLogMessage(publisher.LogMessage{
		Name: "transaction",
		Data: fmt.Sprintf("Transaction success to %s with amount %d", newTransaction.CreditAccount, newTransaction.Amount),
	})

	// send message for email notification
	go ts.publisher.PublishMailMessage(publisher.MailMessage{
		Recipient: "husnir2005@gmail.com",
		DebitAccount: newTransaction.DebitAccount,
		CreditAccount: newTransaction.CreditAccount,
		Amount: int64(newTransaction.Amount),
	})

	return ts.writeResponse(false, "transaction successfully")
}