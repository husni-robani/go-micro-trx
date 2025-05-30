package main

import (
	"encoding/json"
	"log"
	"logger-service/data"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct{
	App Config
	AMQPConn *amqp.Connection
	Tag string
	Queue string
}

type LogData struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func NewConsumer(config Config, amqpConn *amqp.Connection, tag string, queue string) Consumer {
	return Consumer{
		App: config,
		AMQPConn: amqpConn,
		Tag: tag,
		Queue: queue,
	}
}

func (c Consumer) Listen() {
	log.Printf("Start listening to queue:  %s....", c.Queue)

	ch, err := c.AMQPConn.Channel()
	if err != nil {
		log.Fatal("failed to open channel: ", err)
	}

	_, err = ch.QueueDeclare(c.Queue, false, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare queue: ", err)
	}

	msgCh, err := ch.Consume(c.Queue, c.Tag, false, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to consume: ", err)
	}

	
	for msg := range msgCh {
		var newLog LogData

		// unmarshal message
		if err := json.Unmarshal(msg.Body, &newLog); err != nil {
			log.Println("Invalid message body format: ", err)
			msg.Nack(false, false)
			continue
		}

		log.Println("message received: ", newLog)

		if err := c.writeLog(newLog); err != nil {
			log.Println("Failed to write log: ", err)
			msg.Nack(false, false)
			continue
		}

		msg.Ack(false)
	}
}

// write log
func (c Consumer) writeLog(newLog LogData) error {
	log := data.LogEntry{
		Name: newLog.Name,
		Data: newLog.Data,
	}

	if err := c.App.Models.Log.InsertOne(log); err != nil {
		return err
	}
	
	return nil
}