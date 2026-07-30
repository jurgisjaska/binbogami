package auth

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/jurgisjaska/binbogami/internal/database/user/invitation"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockSignupUserRepository struct {
	existingUsers map[string]*user.User
	createdUsers  []*user.User
	failOnCreate  bool
}

func (m *mockSignupUserRepository) FindByEmail(e string) (*user.User, error) {
	if u, ok := m.existingUsers[e]; ok {
		return u, nil
	}
	return nil, nil
}

func (m *mockSignupUserRepository) Create(u *user.User) error {
	if m.failOnCreate {
		return errors.New("database error creating user")
	}
	m.createdUsers = append(m.createdUsers, u)
	return nil
}

func (m *mockSignupUserRepository) Find(id uuid.UUID) (*user.User, error)          { return nil, nil }
func (m *mockSignupUserRepository) FindActive(id uuid.UUID) (*user.User, error)    { return nil, nil }
func (m *mockSignupUserRepository) FindActiveByEmail(e string) (*user.User, error) { return nil, nil }
func (m *mockSignupUserRepository) FindMany(filter string) (*user.Users, error)    { return nil, nil }
func (m *mockSignupUserRepository) UpdatePassword(u *user.User) error              { return nil }

type mockSignupInvitationRepository struct {
	invitations map[uuid.UUID]*invitation.Invitation
	deleted     map[uuid.UUID]*invitation.Invitation
}

func (m *mockSignupInvitationRepository) Open(id uuid.UUID) (*invitation.Invitation, error) {
	inv, ok := m.invitations[id]
	if !ok {
		return nil, errors.New("invitation not found")
	}
	now := time.Now()
	inv.OpenedAt = &now
	return inv, nil
}

func (m *mockSignupInvitationRepository) Find(id uuid.UUID) (*invitation.Invitation, error) {
	inv, ok := m.invitations[id]
	if !ok || inv.DeletedAt != nil {
		return nil, errors.New("invitation not found")
	}
	return inv, nil
}

func (m *mockSignupInvitationRepository) Create(model *models.InvitationRequest) (invitation.Invitations, error) {
	return nil, nil
}

func (m *mockSignupInvitationRepository) Update(inv *invitation.Invitation) error {
	if m.invitations != nil && inv.Id != nil {
		m.invitations[*inv.Id] = inv
	}
	return nil
}

func TestSignup(t *testing.T) {
	e := echo.New()
	e.Validator = &api.Validator{Validator: validator.New()}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	secret := "test-secret-key-32-chars-long-12345"
	config := &internal.Config{
		Secret: secret,
	}

	jackID := uuid.MustParse("05e7257a-b21c-11ee-9a7a-5ab75f0c1cab")
	janetUserID := uuid.MustParse("7c8f9a0b-1234-4567-8901-abcdef123456")

	hammondInvID := uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	harrimanInvID := uuid.MustParse("0685f091-63c3-4d24-8765-e35680621101")
	bratacInvID := uuid.MustParse("b1fec990-1111-4ef8-bb6d-6bb9bd380a22")
	kawalskyInvID := uuid.MustParse("e4cbf223-4444-4ef8-bb6d-6bb9bd380a55")

	roleEditor := 2
	roleAdmin := 4
	now := time.Now()
	exp2028 := time.Date(2028, 1, 1, 1, 1, 1, 0, time.UTC)

	// Existing user fixtures matching database/fixtures.sql
	existingUsers := map[string]*user.User{
		"tealc.of.chulak@sgc.example.com": {
			Id:    uuid.MustParse("aff84550-b21f-11ee-8ac0-5ab75f0c1cab"),
			Email: "tealc.of.chulak@sgc.example.com",
		},
		"janet.fraiser@sgc.example.com": {
			Id:    janetUserID,
			Email: "janet.fraiser@sgc.example.com",
		},
	}

	// Invitation fixtures matching database/fixtures.sql
	invitationsMap := map[uuid.UUID]*invitation.Invitation{
		hammondInvID: {
			Id:        &hammondInvID,
			Email:     "george.hammond@sgc.example.com",
			Role:      nil,
			CreatedBy: &jackID,
			UserId:    nil,
			CreatedAt: now,
			ExpiredAt: exp2028,
		},
		harrimanInvID: {
			Id:        &harrimanInvID,
			Email:     "walter.harriman@sgc.example.com",
			Role:      &roleEditor,
			CreatedBy: &jackID,
			UserId:    nil,
			CreatedAt: now,
			ExpiredAt: exp2028,
		},
		bratacInvID: {
			Id:        &bratacInvID,
			Email:     "bratac@sgc.example.com",
			Role:      &roleEditor,
			CreatedBy: &jackID,
			UserId:    nil,
			CreatedAt: now,
			OpenedAt:  &now,
			ExpiredAt: exp2028,
		},
		kawalskyInvID: {
			Id:        &kawalskyInvID,
			Email:     "charles.kawalsky@sgc.example.com",
			Role:      &roleAdmin,
			CreatedBy: &jackID,
			UserId:    nil,
			CreatedAt: now,
			DeletedAt: &now, // Deleted fixture
			ExpiredAt: exp2028,
		},
	}

	tests := []struct {
		name                 string
		contentType          string
		payload              string
		failUserCreate       bool
		expectedStatus       int
		expectInBody         []string
		checkUserConfirmed   bool
		expectedAssignedRole int
		checkInvitationUser  *uuid.UUID
	}{
		{
			name:           "Invalid JSON payload",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{invalid json`,
			expectedStatus: http.StatusBadRequest,
			expectInBody:   []string{"bad request"},
		},
		{
			name:           "Empty request body",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectInBody:   []string{"error encountered during data validation"},
		},
		{
			name:           "Invalid email format",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"invalid-email","password":"password123","repeated_password":"password123","name":"George","surname":"Hammond"}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectInBody:   []string{"error encountered during data validation"},
		},
		{
			name:           "Short password (< 8 chars)",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"george.hammond@sgc.example.com","password":"123","repeated_password":"123","name":"George","surname":"Hammond"}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectInBody:   []string{"error encountered during data validation"},
		},
		{
			name:           "Mismatched repeated password",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"george.hammond@sgc.example.com","password":"password123","repeated_password":"different123","name":"George","surname":"Hammond"}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectInBody:   []string{"error encountered during data validation"},
		},
		{
			name:           "Existing email fixture (tealc.of.chulak@sgc.example.com)",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"tealc.of.chulak@sgc.example.com","password":"password123","repeated_password":"password123","name":"Teal'c","surname":"Chulak"}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectInBody:   []string{"email address already in use"},
		},
		{
			name:           "Existing email fixture (janet.fraiser@sgc.example.com)",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"janet.fraiser@sgc.example.com","password":"password123","repeated_password":"password123","name":"Janet","surname":"Fraiser"}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectInBody:   []string{"email address already in use"},
		},
		{
			name:           "Invalid invitation UUID format",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"george.hammond@sgc.example.com","password":"password123","repeated_password":"password123","name":"George","surname":"Hammond","invitation_id":"not-a-uuid"}`,
			expectedStatus: http.StatusBadRequest,
			expectInBody:   []string{"bad request"},
		},
		{
			name:           "Revoked/Deleted invitation fixture (charles.kawalsky)",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"charles.kawalsky@sgc.example.com","password":"password123","repeated_password":"password123","name":"Charles","surname":"Kawalsky","invitation_id":"e4cbf223-4444-4ef8-bb6d-6bb9bd380a55"}`,
			expectedStatus: http.StatusOK,
			expectInBody:   []string{"success", "charles.kawalsky@sgc.example.com"},
		},
		{
			name:           "Database error on user creation",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"new.user@sgc.example.com","password":"password123","repeated_password":"password123","name":"New","surname":"User"}`,
			failUserCreate: true,
			expectedStatus: http.StatusInternalServerError,
			expectInBody:   []string{"internal server error"},
		},
		{
			name:                 "Successful signup without invitation",
			contentType:          echo.MIMEApplicationJSON,
			payload:              `{"email":"ronon.dex@sgc.example.com","password":"password123","repeated_password":"password123","name":"Ronon","surname":"Dex"}`,
			expectedStatus:       http.StatusOK,
			expectInBody:         []string{"success", "ronon.dex@sgc.example.com"},
			expectedAssignedRole: user.RoleDefault,
		},
		{
			name:                 "Successful signup with valid invitation fixture (george.hammond)",
			contentType:          echo.MIMEApplicationJSON,
			payload:              `{"email":"george.hammond@sgc.example.com","password":"password123","repeated_password":"password123","name":"George","surname":"Hammond","invitation_id":"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"}`,
			expectedStatus:       http.StatusOK,
			expectInBody:         []string{"success", "george.hammond@sgc.example.com"},
			checkUserConfirmed:   true,
			expectedAssignedRole: user.RoleDefault,
			checkInvitationUser:  invitationsMap[hammondInvID].UserId,
		},
		{
			name:                 "Successful signup with valid invitation fixture with role set (walter.harriman)",
			contentType:          echo.MIMEApplicationJSON,
			payload:              `{"email":"walter.harriman@sgc.example.com","password":"password123","repeated_password":"password123","name":"Walter","surname":"Harriman","invitation_id":"0685f091-63c3-4d24-8765-e35680621101"}`,
			expectedStatus:       http.StatusOK,
			expectInBody:         []string{"success", "walter.harriman@sgc.example.com"},
			checkUserConfirmed:   true,
			expectedAssignedRole: roleEditor,
			checkInvitationUser:  invitationsMap[harrimanInvID].UserId,
		},
		{
			name:                 "Successful signup with valid opened invitation fixture (bratac)",
			contentType:          echo.MIMEApplicationJSON,
			payload:              `{"email":"bratac@sgc.example.com","password":"password123","repeated_password":"password123","name":"Master","surname":"Bra'tac","invitation_id":"b1fec990-1111-4ef8-bb6d-6bb9bd380a22"}`,
			expectedStatus:       http.StatusOK,
			expectInBody:         []string{"success", "bratac@sgc.example.com"},
			checkUserConfirmed:   true,
			expectedAssignedRole: roleEditor,
			checkInvitationUser:  invitationsMap[bratacInvID].UserId,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &mockSignupUserRepository{
				existingUsers: existingUsers,
				failOnCreate:  tt.failUserCreate,
			}
			mockInvRepo := &mockSignupInvitationRepository{
				invitations: invitationsMap,
			}

			h := &Auth{
				echo:          e,
				configuration: config,
				auditlog:      logger,
				invitation:    mockInvRepo,
				user: &userRepositories{
					repository: mockUserRepo,
				},
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBufferString(tt.payload))
			if tt.contentType != "" {
				req.Header.Set(echo.HeaderContentType, tt.contentType)
			}
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			err := h.signup(c)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			for _, expectedStr := range tt.expectInBody {
				assert.Contains(t, rec.Body.String(), expectedStr)
			}

			if tt.expectedStatus == http.StatusOK && len(mockUserRepo.createdUsers) > 0 {
				createdUser := mockUserRepo.createdUsers[len(mockUserRepo.createdUsers)-1]
				if tt.checkUserConfirmed {
					assert.NotNil(t, createdUser.ConfirmedAt)
				}
				if tt.expectedAssignedRole != 0 {
					assert.Equal(t, tt.expectedAssignedRole, createdUser.Role)
				}
				if tt.checkInvitationUser != nil {
					assert.Equal(t, createdUser.Id, *tt.checkInvitationUser)
				}
			}
		})
	}
}
