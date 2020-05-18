package email

import (
	"github.com/matcornic/hermes/v2"
	"gopkg.in/gomail.v2"
	"log"
)

type Mail struct {
	Subject string
	Body    string
	To      string
}

type Product struct {
	Name string
	Link string
	Logo string
}

type Button struct {
	Color string
	Text  string
	Link  string
}

type Action struct {
	Instructions string
	Button       Button
}

type Entry struct {
	Key   string
	Value string
}

type MailData struct {
	Name    string
	Intros  []string
	Actions []Action
	Outros  []string
	Header  Product
	Entry   [][]Entry
}

func (b MailData) Generate() (string, error) {

	header := hermes.Hermes{
		Product: hermes.Product{
			Name:        "Afiqo",
			Copyright:   "Copyright © 2020 Afiqo-Location. All rights reserved.",
			Logo:        "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
			TroubleText: "Feel free to contact us at +60122872072",
		},
	}

	var action []hermes.Action

	for _, ac := range b.Actions {
		action = append(action, hermes.Action{
			Instructions: ac.Instructions,
			Button: hermes.Button{
				Color: ac.Button.Color,
				Text:  ac.Button.Text,
				Link:  ac.Button.Link,
			},
		})
	}

	emailTemplate := hermes.Email{
		Body: hermes.Body{
			Name:    b.Name,
			Intros:  b.Intros,
			Actions: action,
			Outros:  b.Outros,
		},
	}

	emailBody, err := header.GenerateHTML(emailTemplate)
	if err != nil {
		return "", err
	}

	return emailBody, err
}

func (b MailData) GenerateForPassword() (string, error) {

	header := hermes.Hermes{
		Product: hermes.Product{
			Name:        "Afiqo",
			Copyright:   "Copyright © 2020 Afiqo-Location. All rights reserved.",
			Logo:        "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
			TroubleText: "Feel free to contact us at +60123456789",
		},
	}

	var action []hermes.Action

	for _, ac := range b.Actions {
		action = append(action, hermes.Action{
			Instructions: "The password for your login is stated below:",
			Button: hermes.Button{
				Color: "#2732b0",
				Text:  ac.Button.Text,
			},
		})
	}

	emailTemplate := hermes.Email{
		Body: hermes.Body{
			Name: b.Name,
			Intros: []string{
				"Welcome to Afiqo-Location! We're very excited to have you on board.",
			},
			Actions: action,
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}

	emailBody, err := header.GenerateHTML(emailTemplate)
	if err != nil {
		return "", err
	}

	return emailBody, err
}

func (b MailData) GenerateForReceipt() (string, error) {

	header := hermes.Hermes{
		Product: hermes.Product{
			Name:        "Afiqo",
			Copyright:   "Copyright © 2020 Afiqo-Location. All rights reserved.",
			Logo:        "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
			TroubleText: "Feel free to contact us at +60123456789",
		},
	}

	var entries [][]hermes.Entry

	for _, entry := range b.Entry {
		var rows []hermes.Entry
		for _, value := range entry {
			row := hermes.Entry{
				Key:   value.Key,
				Value: value.Value,
			}
			rows = append(rows, row)
		}
		entries = append(entries, rows)
	}

	emailTemplate := hermes.Email{
		Body: hermes.Body{
			Name: b.Name,
			Intros: []string{
				"Your order is being processed.",
			},

			Table: hermes.Table{
				Data: entries,
				Columns: hermes.Columns{
					CustomWidth: map[string]string{
						"Item":     "15%",
						"Price":    "15%",
						"Subtotal": "15%",
						"Quantity": "10%",
					},
					CustomAlignment: map[string]string{},
				},
			},

			Actions: []hermes.Action{
				{
					Instructions: "You can check the status of your order and more in your dashboard:",
					Button: hermes.Button{
						Text: "Go to Dashboard",
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}

	emailBody, err := header.GenerateHTML(emailTemplate)
	if err != nil {
		return "", err
	}

	return emailBody, err
}

func (m Mail) SendEmail() {

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", email)
	mailer.SetHeader("To", m.To)
	mailer.SetHeader("Subject", m.Subject)
	mailer.SetBody("text/html", m.Body)

	dialer := gomail.NewDialer(
		host,
		port,
		email,
		password,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Mail sent!")
}
