package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"runtime"

	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/jurgisjaska/binbogami/internal/database/user/invitation"
	"gopkg.in/gomail.v2"
)

type (
	Invitation struct {
		d *gomail.Dialer
		c *internal.Config
	}

	InvitationContent struct {
		URL          string
		Sender       string
		Organization string
	}
)

func (m *Invitation) Send(sender *user.User, i *invitation.Invitation) error {
	message := gomail.NewMessage()
	message.SetHeader("From", m.c.Mail.Sender)
	message.SetHeader("To", i.Email)
	message.SetHeader("Subject", "Invitation")

	content, err := m.createMessage(sender, i)
	if err != nil {
		return err
	}

	message.SetBody("text/html", content)

	return m.d.DialAndSend(message)
}

func (m *Invitation) createMessage(sender *user.User, i *invitation.Invitation) (string, error) {
	_, f, _, _ := runtime.Caller(0)
	dir := filepath.Dir(f)

	t, err := template.ParseFiles(filepath.Join(dir, "../../../var/templates/invitation.html"))
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("http://%s:%d/signup/%s", m.c.Web.Hostname, m.c.Web.Port, i.Id.String())
	content := InvitationContent{
		URL:    url,
		Sender: fmt.Sprintf("%s %s", sender.Name, sender.Surname),
	}

	var b bytes.Buffer
	err = t.Execute(&b, content)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func CreateInvitation(d *gomail.Dialer, c *internal.Config) *Invitation {
	return &Invitation{d: d, c: c}
}
