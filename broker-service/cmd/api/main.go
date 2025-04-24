package main

import (
	"log"
	"net/http"
)

const(
	webPort = "80"
)

type Config struct{}

func main() {
	log.Println("Running broker-service on port ", webPort)

	app := Config{}

	srv := &http.Server{
		Addr: ":" + webPort,
		Handler: app.routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic("Failed to run broker-service")
	}
}