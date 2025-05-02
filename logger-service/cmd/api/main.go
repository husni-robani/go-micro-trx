package main

import (
	"log"
	"logger-service/data"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	mongo_uri = "mongodb://mongo:27017"
	web_port = "80"
)

type Config struct {
	Models data.Models
}

func main() {
	mongoClient := connectMongoDB()
	app := Config{
		Models: data.New(mongoClient),
	}

	log.Println("Starting Web Server ...")

	srv := http.Server{
		Addr: ":" + web_port,
		Handler: app.routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic("Failed to running web server: ", err)
	}
}

func connectMongoDB() *mongo.Client {
	log.Println("Starting MongoDB Connection ...")

	client, err := mongo.Connect(options.Client().ApplyURI(mongo_uri))
	if err != nil {
		log.Println("Failed to connect MongoDB: ", err)
	}

	return client
}