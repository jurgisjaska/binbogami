package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"runtime"

	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"gopkg.in/gomail.v2"
)

type ResetPassword struct {
	d *gomail.Dialer
	c *internal.Config
}

func (m *ResetPassword) Send(u *user.User, pr *user.PasswordReset) error {
	message := gomail.NewMessage()
	message.SetHeader("From", m.c.Mail.Sender)
	message.SetHeader("To", *u.Email)
	message.SetHeader("Subject", "Reset Password")

	content, err := m.createMessage(*u.Name, pr.Id)
	if err != nil {
		return err
	}

	message.SetBody("text/html", content)

	return m.d.DialAndSend(message)
}

func (m *ResetPassword) createMessage(u string, id *uuid.UUID) (string, error) {
	_, f, _, _ := runtime.Caller(0)
	dir := filepath.Dir(f)

	t, err := template.ParseFiles(filepath.Join(dir, "../../../var/template/reset_password.html"))
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("http://%s:%d/reset-password/%s", m.c.Hostname, m.c.Port, id.String())
	content := struct {
		Name string
		URL  string
	}{
		Name: u,
		URL:  url,
	}

	var b bytes.Buffer
	err = t.Execute(&b, content)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

// CreateResetPassword creates a new instance of the ResetPassword mail service.
func CreateResetPassword(d *gomail.Dialer, c *internal.Config) *ResetPassword {
	return &ResetPassword{d: d, c: c}
}
