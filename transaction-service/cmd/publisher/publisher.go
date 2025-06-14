package publisher

import (
	"encoding/json"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	AMQPConn *amqp.Connection
	Exchange string
}

type LogMessage struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailMessage struct {
	Recipient string `json:"recipient"`
	DebitAccount string `json:"debit_account"`
	CreditAccount string `json:"credit_account"`
	Amount int64 `json:"amount"`
}

func NewPublisher(amqpConn *amqp.Connection, exchange string) Publisher {
	return Publisher{
		AMQPConn: amqpConn,
		Exchange: exchange,
	}
}

func (p Publisher) PublishMailMessage(m MailMessage) error {
	routingKey := os.Getenv("MAIL_ROUTING_KEY")

	messageJson, err := json.Marshal(m)
	if err != nil {
		log.Println("Failed to marshal message: ", err)
		return err
	}

	return p.publish(routingKey, messageJson)
}

func (p Publisher) PublishLogMessage(l LogMessage) error {
	routingKey := os.Getenv("LOG_ROUTING_KEY")

	messageJson, err := json.Marshal(l)
	if err != nil {
		log.Println("Failed to marshal message: ", err)
		return err
	}

	return p.publish(routingKey, messageJson)
}

func (p Publisher) publish(routingKey string, message []byte) error {
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

	err = ch.Publish(p.Exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body: message,
	})
	
	if err != nil {
		log.Println("Failed to publish message: ", err)
		return err
	}

	return nil	
}

// func (p Publisher) PublishMessage(jsonMessage []byte) error {
// 	ch, err := p.AMQPConn.Channel()
// 	if err != nil {
// 		log.Println("Failed to open channel: ", err)
// 		return err
// 	}

// 	err = ch.ExchangeDeclare(p.Exchange, "direct", true, false, false, false, nil)
// 	if err != nil {
// 		log.Println("Failed to declare exchange: ", err)
// 		return err
// 	}

// 	err = ch.Publish(p.Exchange, p.RoutingKey, false, false, amqp.Publishing{
// 		ContentType: "application/json",
// 		Body: jsonMessage,
// 	})
	
// 	if err != nil {
// 		log.Println("Failed to publish message: ", err)
// 		return err
// 	}

// 	return nil
// }