package main

import "transaction-service/data"

type TransactionRPCServer struct {
	models data.Models
}

func NewTransactionRPCServer(m data.Models) TransactionRPCServer {
	return TransactionRPCServer{
		models: m,
	}
}