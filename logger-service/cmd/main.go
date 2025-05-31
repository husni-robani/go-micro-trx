package main

import (
	"log"
	"logger-service/data"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Config struct {
	Models data.Models
}

func main() {
	mongoClient := connectMongoDB()
	
	app := Config{
		Models: data.New(mongoClient),
	}

	consumer := NewConsumer(app, connectRabbitMQ(), "worker-1", os.Getenv("QUEUE_NAME"), os.Getenv("EXCHANGE_NAME"), os.Getenv("ROUTING_KEY"))
	consumer.Listen()
}

func connectRabbitMQ() *amqp.Connection {
	attemp := 0

	for {
		conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
		if err != nil {
			if attemp >= 5 {
				log.Fatal("RabbitMQ connection failed: ", err)
				return nil
			}

			attemp ++
			log.Println("AMQP connection not ready yet....")
			time.Sleep(time.Second * 3)
			continue
		}

		return conn
	}
}

func connectMongoDB() *mongo.Client {
	log.Println("Starting MongoDB Connection ...")

	option := options.Client().ApplyURI(os.Getenv("MONGO_DSN"))
	option.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	client, err := mongo.Connect(option)
	if err != nil {
		log.Panic("Failed to connect MongoDB: ", err)
	}

	log.Println("MongoDB connected!")

	return client
}