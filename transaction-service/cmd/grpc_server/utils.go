package grpc_server

import (
	transactionpb "proto/transaction"
)

func (ts TransactionGRPCServer) writeErrorResponse(err error) (*transactionpb.TransactionResponse, error) {

	return nil, err
}

func (ts TransactionGRPCServer) writeResponse(isError bool, message string) (*transactionpb.TransactionResponse, error){
	var res transactionpb.TransactionResponse
	res.Error = isError
	res.Message = message
	return &res, nil
}