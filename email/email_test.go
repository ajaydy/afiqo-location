package email

import (
	"gopkg.in/gomail.v2"
	"log"
	"testing"
)

//remove password

func TestSendEmail(t *testing.T) {

	data := MailData{
		Name: "Afiq",
		Entry: [][]Entry{
			{
				{
					Key: "Item", Value: "Golang",
				},
				{
					Key: "Description", Value: "Open Source Programming",
				},
				{
					Key: "Price", Value: "13",
				},
			},
			{
				{
					Key: "", Value: "Hermes",
				},
				{
					Key: "Description", Value: "Programmatically create beautiful e-mails using Golang.",
				},
				{
					Key: "Price", Value: "$1.99",
				},
			},
		},
	}

	body, err := data.GenerateForReceipt()

	if err != nil {
		log.Fatal(err)
	}

	mail := MailService{
		Host:     "smtp.gmail.com",
		Port:     587,
		Email:    "hannahmohamed220800@gmail.com",
		Password: "",
	}

	mail.Init()

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", email)
	mailer.SetHeader("To", "mohdjamilafiq@gmail.com")
	mailer.SetHeader("Subject", "Test mail")
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(
		host,
		port,
		email,
		password,
	)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Mail sent!")
}
