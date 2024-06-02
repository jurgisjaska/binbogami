package internal

import (
	"fmt"
	"net/smtp"

	"gopkg.in/gomail.v2"
)

// @deprecated use dialer and gomail package with mailing services.
func ConnectMail(c *Mail) (*smtp.Client, error) {
	dsn := fmt.Sprintf("%s:%d", c.Connection.Hostname, c.Connection.Port)
	client, err := smtp.Dial(dsn)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func CreateDialer(c *Mail) *gomail.Dialer {
	return gomail.NewDialer(c.Connection.Hostname, c.Connection.Port, c.Connection.Username, c.Connection.Password)
}
