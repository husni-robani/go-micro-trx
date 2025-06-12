package config

import (
	"net/rpc"
	taskpb "proto/task"
)

type Config struct {
	RPCClientTask *rpc.Client
	GRPCClientTask taskpb.TaskServiceClient
}