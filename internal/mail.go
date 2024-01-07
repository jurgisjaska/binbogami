package internal

import (
	"fmt"
	"net/smtp"
)

func ConnectMail(c *Mail) (*smtp.Client, error) {
	dsn := fmt.Sprintf("%s:%d", c.Connection.Hostname, c.Connection.Port)
	client, err := smtp.Dial(dsn)
	if err != nil {
		return nil, err
	}

	return client, nil
}
