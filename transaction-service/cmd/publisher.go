package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	AMQPConn *amqp.Connection
	Exchange string
	RoutingKey string
}

type LogPayload struct {
	Name string
	Data string
}

func NewPublisher(amqpConn *amqp.Connection, exchange string, routingKey string) Publisher {
	return Publisher{
		AMQPConn: amqpConn,
		Exchange: exchange,
		RoutingKey: routingKey,
	}
}

func (p Publisher) PublishMessage(jsonMessage []byte) error {
	ch, err := p.AMQPConn.Channel()
	if err != nil {
		log.Println("Failed to open channel: ", err)
		return err
	}

	// q, err := ch.QueueDeclare(os.Getenv("LOG_QUEUE_NAME"), false, false, false, false, nil)
	// if err != nil {
	// 	log.Println("Failed to declare queue: ", err)
	// 	return err
	// }

	err = ch.Publish("", p.RoutingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body: jsonMessage,
	})
	
	if err != nil {
		log.Println("Failed to publish message: ", err)
		return err
	}

	return nil
}