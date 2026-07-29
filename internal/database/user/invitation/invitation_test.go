package invitation

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInvitationStruct(t *testing.T) {
	id := uuid.New()
	createdBy := uuid.MustParse("05e7257a-b21c-11ee-9a7a-5ab75f0c1cab") // Jack O'Neill
	userId := uuid.MustParse("7c8f9a0b-1234-4567-8901-abcdef123456")    // Janet Fraiser
	role := 4
	now := time.Now()
	expired := time.Date(2028, 1, 1, 1, 1, 1, 0, time.UTC)

	tests := []struct {
		name       string
		invitation *Invitation
	}{
		{
			name: "Pending Invitation with null role and user",
			invitation: &Invitation{
				Id:        &id,
				Email:     "george.hammond@sgc.example.com",
				Role:      nil,
				CreatedBy: &createdBy,
				UserId:    nil,
				CreatedAt: now,
				OpenedAt:  nil,
				DeletedAt: nil,
				ExpiredAt: expired,
			},
		},
		{
			name: "Accepted Invitation with role and accepted user ID",
			invitation: &Invitation{
				Id:        &id,
				Email:     "janet.fraiser@sgc.example.com",
				Role:      &role,
				CreatedBy: &createdBy,
				UserId:    &userId,
				CreatedAt: now,
				OpenedAt:  &now,
				DeletedAt: &now,
				ExpiredAt: expired,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.invitation.Id)
			assert.NotEmpty(t, tt.invitation.Email)
			assert.Equal(t, &createdBy, tt.invitation.CreatedBy)
			assert.True(t, tt.invitation.ExpiredAt.Year() >= 2028)
		})
	}
}
