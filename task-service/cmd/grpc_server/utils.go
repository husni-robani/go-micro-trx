package grpc_server

import (
	taskpb "proto/task"
)

func (ts TaskGRPCServer) writeErrorResponse(err error) (*taskpb.TaskResponse, error) {

	return nil, err
}

func (ts TaskGRPCServer) writeResponse(isError bool, message string) (*taskpb.TaskResponse, error){
	var res taskpb.TaskResponse
	res.Error = isError
	res.Message = message
	return &res, nil
}