package email

import (
	"encoding/base64"
	"fmt"
)

type MailService struct {
	Host     string
	Port     int
	Email    string
	Password string
}

var (
	host     string
	port     int
	email    string
	password string
)

func (mail MailService) Init() {
	decoded, err := base64.StdEncoding.DecodeString(mail.Password)
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}
	host = mail.Host
	port = mail.Port
	email = mail.Email
	password = string(decoded)
}
