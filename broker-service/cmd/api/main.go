package main

import (
	"log"
	"net/http"
	"net/rpc"
	"time"
)

const(
	webPort = "80"
)

type Config struct{
	rpcClient *rpc.Client
}

func main() {
	log.Println("Running broker-service on port ", webPort)

	app := Config{
		rpcClient: connectRPC(),
	}

	srv := &http.Server{
		Addr: ":" + webPort,
		Handler: app.routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic("Failed to run broker-service")
	}
}

func connectRPC() *rpc.Client{
	counter := 1

	for {
		client, err := rpc.Dial("tcp", "task-service:50001")
		if err != nil {
			if counter >= 5 {
				log.Panic("RPC Connection failed: ", err)
			}
			log.Printf("RPC not connected yet: %v (waiting for 2 second)\n", err)
			time.Sleep(time.Second * 2)
			counter ++
			continue
		}

		log.Println("RPC Connected!")
		
		return client
	}

}