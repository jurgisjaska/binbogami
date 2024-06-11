package internal

import (
	"gopkg.in/gomail.v2"
)

func CreateDialer(c *Mail) *gomail.Dialer {
	return gomail.NewDialer(c.Connection.Hostname, c.Connection.Port, c.Connection.Username, c.Connection.Password)
}
