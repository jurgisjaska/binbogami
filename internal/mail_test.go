package internal

import (
	"testing"

	"gopkg.in/gomail.v2"
)

type mockDialer struct{}

func (md *mockDialer) DialAndSend(mm ...*gomail.Message) error {
	return nil
}

func TestCreateDialer(t *testing.T) {
	tableTests := []struct {
		name   string
		input  *Mail
		expect *gomail.Dialer
	}{
		{
			"Test valid mail",
			&Mail{
				"info@example.com",
				&Connection{
					Hostname: "smtp.mail.com",
					Port:     587,
					Username: "username@mail.com",
					Password: "password",
				},
			},
			&gomail.Dialer{
				Host:     "smtp.mail.com",
				Port:     587,
				Username: "username@mail.com",
				Password: "password",
			},
		},
		{
			"Test empty mail",
			&Mail{
				"",
				&Connection{
					Hostname: "",
					Port:     0,
					Username: "",
					Password: "",
				},
			},
			&gomail.Dialer{
				Host:     "",
				Port:     0,
				Username: "",
				Password: "",
			},
		},
	}

	for _, tt := range tableTests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateDialer(tt.input)
			if got.Host != tt.expect.Host || got.Port != tt.expect.Port || got.Username != tt.expect.Username || got.Password != tt.expect.Password {
				t.Errorf("CreateDialer() = %v, want %v", got, tt.expect)
			}
		})
	}
}
