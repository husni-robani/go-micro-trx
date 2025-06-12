package main

import (
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	amqpConn := connectRabbitMQ()

	consumer := NewConsumer(amqpConn, "worker-1", os.Getenv("QUEUE_NAME"), os.Getenv("EXCHANGE_NAME"), os.Getenv("ROUTING_KEY"))

	consumer.Listen()
}

func connectRabbitMQ() *amqp.Connection {
	counter := 0

	for {
		conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
		if err != nil {
			counter ++

			if counter >= 5 {
				log.Fatal("RabbitMQ connection failed: ", err)
			}

			log.Println("RabbitMQ not connected yet ...")
			time.Sleep(time.Second * 3)
			continue
		}

		return conn
	}
}