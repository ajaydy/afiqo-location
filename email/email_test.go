package email

import (
	"gopkg.in/gomail.v2"
	"log"
	"testing"
)

func TestSendEmail(t *testing.T) {

	data := MailData{
		Name: "Afiq Jamil",

		Actions: []Action{
			{
				Button: Button{
					Text: "1626273883",
				},
			},
		},
	}

	body, err := data.GenerateForPassword()

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
