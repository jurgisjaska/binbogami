package auth

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/database/user/invitation"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockInvitationRepository struct {
	invitationsByID map[uuid.UUID]*invitation.Invitation
	failOnUpdate    bool
}

func (m *mockInvitationRepository) Find(id uuid.UUID) (*invitation.Invitation, error) {
	inv, ok := m.invitationsByID[id]
	if !ok {
		return nil, errors.New("invitation not found")
	}
	return inv, nil
}

func (m *mockInvitationRepository) Update(i *invitation.Invitation) error {
	if m.failOnUpdate {
		return errors.New("failed to update invitation")
	}
	if m.invitationsByID != nil && i.Id != nil {
		m.invitationsByID[*i.Id] = i
	}
	return nil
}

func TestOpenInvitation(t *testing.T) {
	e := echo.New()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	secret := "test-secret-key-32-chars-long-12345"
	config := &internal.Config{
		Secret: secret,
	}

	validID := uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	openedID := uuid.MustParse("b1fec990-1111-4ef8-bb6d-6bb9bd380a22")

	now := time.Now()
	unopenedInvitation := &invitation.Invitation{
		Id:        &validID,
		Email:     "george.hammond@sgc.example.com",
		CreatedAt: now,
		ExpiredAt: now.Add(24 * time.Hour),
		OpenedAt:  nil,
	}

	openedInvitation := &invitation.Invitation{
		Id:        &openedID,
		Email:     "bratac@sgc.example.com",
		CreatedAt: now,
		ExpiredAt: now.Add(24 * time.Hour),
		OpenedAt:  &now,
	}

	tests := []struct {
		name           string
		paramID        string
		failOnUpdate   bool
		expectedStatus int
		expectInBody   []string
		checkOpenedAt  bool
	}{
		{
			name:           "Invalid invitation UUID format",
			paramID:        "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
			expectInBody:   []string{"incorrect invitation"},
		},
		{
			name:           "Invitation not found",
			paramID:        "11111111-1111-1111-1111-111111111111",
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"invitation not found"},
		},
		{
			name:           "Database error updating invitation",
			paramID:        validID.String(),
			failOnUpdate:   true,
			expectedStatus: http.StatusInternalServerError,
			expectInBody:   []string{"failed to update invitation"},
		},
		{
			name:           "Successful open of unopened invitation",
			paramID:        validID.String(),
			expectedStatus: http.StatusOK,
			expectInBody:   []string{"success", validID.String(), "george.hammond@sgc.example.com"},
			checkOpenedAt:  true,
		},
		{
			name:           "Successful open of already opened invitation",
			paramID:        openedID.String(),
			expectedStatus: http.StatusOK,
			expectInBody:   []string{"success", openedID.String(), "bratac@sgc.example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unopenedCopy := *unopenedInvitation
			openedCopy := *openedInvitation

			mockRepo := &mockInvitationRepository{
				invitationsByID: map[uuid.UUID]*invitation.Invitation{
					validID:  &unopenedCopy,
					openedID: &openedCopy,
				},
				failOnUpdate: tt.failOnUpdate,
			}

			h := &Auth{
				echo:          e,
				configuration: config,
				auditlog:      logger,
				invitation:    mockRepo,
			}

			req := httptest.NewRequest(http.MethodGet, "/auth/invitation/"+tt.paramID, nil)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath("/auth/invitation/:id")
			c.SetPathValues(echo.PathValues{echo.PathValue{Name: "id", Value: tt.paramID}})

			err := h.openInvitation(c)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			for _, expectedStr := range tt.expectInBody {
				assert.Contains(t, rec.Body.String(), expectedStr)
			}

			if tt.checkOpenedAt {
				assert.NotNil(t, unopenedCopy.OpenedAt)
			}
		})
	}
}
