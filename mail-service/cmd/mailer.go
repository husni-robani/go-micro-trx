package main

import (
	"log"
	"os"
	"strconv"

	"github.com/go-gomail/gomail"
)

type Mailer struct {
	To string
	Subject string
	Body string
}

func (m Mailer) SendEmailNotification() error {
	host := os.Getenv("SMTP_HOST")
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Println("Failed to convert port: ", err)
		return err
	}

	log.Printf("Sending email to %s ....", m.To)

	message := gomail.NewMessage()
	message.SetHeader("From", os.Getenv("EMAIL_ADDRESS"))
	message.SetHeader("To", m.To)
	message.SetHeader("Subject", "Hello!")
	message.SetBody("plain/text", m.Body)

	dialer := gomail.NewDialer(host, port, os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"))
	if err := dialer.DialAndSend(message); err != nil {
		log.Println("Failed to send email: ", err)
		return err
	}

	log.Println("Email sent to ", m.To)

	return nil
}