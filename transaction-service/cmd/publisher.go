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

	err = ch.ExchangeDeclare(p.Exchange, "direct", true, false, false, false, nil)
	if err != nil {
		log.Println("Failed to declare exchange: ", err)
		return err
	}

	err = ch.Publish(p.Exchange, p.RoutingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body: jsonMessage,
	})
	
	if err != nil {
		log.Println("Failed to publish message: ", err)
		return err
	}

	return nil
}