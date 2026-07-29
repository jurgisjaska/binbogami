package mail

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/jurgisjaska/binbogami/internal/database/user/invitation"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gomail.v2"
)

func TestCreateInvitation(t *testing.T) {
	var dialer *gomail.Dialer
	var config *internal.Config

	result := CreateInvitation(dialer, config)

	assert.Equal(t, dialer, result.d)
	assert.Equal(t, config, result.c)
}

func TestInvitation_CreateMessage(t *testing.T) {
	invID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	tests := []struct {
		name    string
		sender  *user.User
		inv     *invitation.Invitation
		wantSub []string
	}{
		{
			name: "Valid invitation message rendering",
			sender: &user.User{
				Name:    "John",
				Surname: "Doe",
			},
			inv: &invitation.Invitation{
				Id: &invID,
			},
			wantSub: []string{
				"John Doe invited you to join personal finance management tool.",
				"http://localhost:8080/signup/00000000-0000-0000-0000-000000000002",
				"logo.png",
				"You've Been Invited",
				"#206bc4",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Invitation{
				c: &internal.Config{
					Web: &internal.URI{
						Hostname: "localhost",
						Port:     8080,
					},
				},
			}

			msg, err := svc.createMessage(tt.sender, tt.inv)
			assert.NoError(t, err)
			for _, sub := range tt.wantSub {
				assert.Contains(t, msg, sub)
			}
		})
	}
}
