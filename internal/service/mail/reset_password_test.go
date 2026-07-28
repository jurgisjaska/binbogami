package mail

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gomail.v2"
)

func TestCreateResetPassword(t *testing.T) {
	var dialer *gomail.Dialer
	var config *internal.Config

	result := CreateResetPassword(dialer, config)

	assert.Equal(t, dialer, result.d)
	assert.Equal(t, config, result.c)
}

func TestResetPassword_CreateMessage(t *testing.T) {
	tests := []struct {
		name     string
		userName string
		id       uuid.UUID
		wantSub  []string
	}{
		{
			name:     "Valid reset message rendering",
			userName: "Jonas Quinn",
			id:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			wantSub: []string{
				"Dear Jonas Quinn,",
				"http://localhost:8080/reset-password/00000000-0000-0000-0000-000000000001",
				"data:image/png;base64,",
				"Reset Your Password",
				"#206bc4",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &ResetPassword{
				c: &internal.Config{
					Web: &internal.URI{
						Hostname: "localhost",
						Port:     8080,
					},
				},
			}

			msg, err := svc.createMessage(tt.userName, &tt.id)
			assert.NoError(t, err)
			for _, sub := range tt.wantSub {
				assert.Contains(t, msg, sub)
			}
		})
	}
}
