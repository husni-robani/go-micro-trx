package main

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	AMQPConn *amqp.Connection
	Tag string
	Queue string 
	Exchange string
	RoutingKey string
}

type MailMessage struct {
	Recipient string `json:"recipient"`
	DebitAccount string `json:"debit_account"`
	CreditAccount string `json:"credit_account"`
	Amount int64 `json:"amount"`
}

func NewConsumer(amqpConn *amqp.Connection, tag string, queue string, exchange string, routing_key string) Consumer {
	return Consumer{
		AMQPConn: amqpConn,
		Tag: tag,
		Queue: queue,
		Exchange: exchange,
		RoutingKey: routing_key,
	}
}

func (c Consumer) Listen() {
	log.Printf("Start listening to queue %s ....\n", c.Queue)

	ch, err := c.AMQPConn.Channel()
	if err != nil {
		log.Fatal("Open channel failed: ", err)
	}

	// declare exchange
	err = ch.ExchangeDeclare(c.Exchange, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Declare Exchange failed: ", err)
	}

	// declare queue
	q, err := ch.QueueDeclare(c.Queue, false, false, false, false, nil)
	if err != nil {
		log.Fatal("Declare Queue failed: ", err)
	}

	// bindqueue
	err = ch.QueueBind(q.Name, c.RoutingKey, c.Exchange, false, nil)
	if err != nil {
		log.Fatal("Binding Queue failed: ", err)
	}

	// consume
	msgCh, err := ch.Consume(q.Name, c.Tag, true, false, false, false, nil)
	if err != nil {
		log.Fatal("Consume Initiation failed: ", err)
	}

	for msg := range msgCh {
		var message MailMessage
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Println("Failed to unmarshal message: ", err)
			return
		}
		

		// Send email
		mailer := Mailer{To: message.Recipient, Subject: "Transaction Notification", Body: fmt.Sprintf("Transaction success to %s with amount %v", message.CreditAccount, message.Amount)}

		go mailer.SendEmailNotification()
	}

}