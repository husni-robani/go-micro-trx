package main

import (
	"broker-service/cmd/api"
	"broker-service/cmd/config"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"time"
)

func main() {
	app := config.Config{
		RPCClientTask: connectRPCClient(os.Getenv("TASK_RPC_ADDRESS")),
	}

	api := api.APIHandler{
		App: &app,
	}


	serveAPI(api.Routes())
}

func serveAPI(routes http.Handler) {
	port := os.Getenv("HTTP_PORT")

	log.Printf("Starting HTTP on port %v....", port)

	server := http.Server{
		Addr: ":" + port,
		Handler: routes,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("HTTP | Failed to run server on port %v", port)
	}
}

func connectRPCClient(address string) *rpc.Client{
	counter := 1

	for {
		client, err := rpc.Dial("tcp", address)
		if err != nil {
			if counter >= 5 {
				log.Panicf("RPC | %v Connection failed: %v", address, err)
			}
			log.Printf("RPC not connected yet: %v (waiting for 2 second)\n", err)
			time.Sleep(time.Second * 2)
			counter ++
			continue
		}

		log.Printf("%v Connected!", address)
		
		return client
	}

}