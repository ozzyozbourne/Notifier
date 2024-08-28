package main

import (
	"context"
	"fmt"
	"log"
	"net/smtp"

	"github.com/twmb/franz-go/pkg/kgo"
)

func sendEmail(to string, subject string, body string) {
	from := "khanosaid726@gmail.com"
	password := "mmcv jhso lmrz ywrh"

	// Gmail's SMTP server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Create the authentication for the SendMail() function
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Create the email content.
	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	// Send the email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		log.Fatalf("Failed to send email: %s", err)
	}

	fmt.Println("Email sent successfully!")
}

func main() {
	to := "khanosaid726@gmail.com"
	subject := "Test Subject"
	body := "This is the body of the email."

	sendEmail(to, subject, body)
}
