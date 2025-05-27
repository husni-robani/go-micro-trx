package config

import "net/rpc"

type Config struct {
	RPCClientTask *rpc.Client
}