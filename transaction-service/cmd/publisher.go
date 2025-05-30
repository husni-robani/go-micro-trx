package main

import (
	"encoding/json"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	AMQPConn *amqp.Connection
}

type LogPayload struct {
	Name string
	Data string
}

func (p Publisher) PublishLogMessage(name string, data string) error {
	payload := LogPayload{Name: name, Data: data}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to marshal payload: ", err)
		return err
	}

	ch, err := p.AMQPConn.Channel()
	if err != nil {
		log.Println("Failed to open channel: ", err)
		return err
	}

	q, err := ch.QueueDeclare(os.Getenv("LOG_QUEUE_NAME"), false, false, false, false, nil)
	if err != nil {
		log.Println("Failed to declare queue: ", err)
		return err
	}

	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body: jsonPayload,
	})
	if err != nil {
		log.Println("Failed to publish message: ", err)
		return err
	}

	return nil
}