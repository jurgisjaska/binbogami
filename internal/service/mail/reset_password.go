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
	"github.com/jurgisjaska/binbogami/internal/database/user/password"
	"gopkg.in/gomail.v2"
)

// ResetPassword represents a mail service for sending password reset emails.
type ResetPassword struct {
	d *gomail.Dialer
	c *internal.Config
}

func (m *ResetPassword) Send(u *user.User, pr *password.Reset) error {
	message := gomail.NewMessage()
	message.SetHeader("From", m.c.Mail.Sender)
	message.SetHeader("To", u.Email)
	message.SetHeader("Subject", "Reset Password")

	content, err := m.createMessage(u.Name, &pr.Id)
	if err != nil {
		return err
	}

	message.SetBody("text/html", content)

	return m.d.DialAndSend(message)
}

func (m *ResetPassword) createMessage(u string, id *uuid.UUID) (string, error) {
	_, f, _, _ := runtime.Caller(0)
	dir := filepath.Dir(f)

	// trying to remember reason why this was moved ???
	// was it because templates are static?
	t, err := template.ParseFiles(filepath.Join(dir, "../../../templates/reset_password.html"))
	if err != nil {
		return "", err
	}

	// password must be sent to the client apps
	// in this case send it to the web application
	url := fmt.Sprintf("http://%s:%d/reset-password/%s", m.c.Web.Hostname, m.c.Web.Port, id.String())
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
