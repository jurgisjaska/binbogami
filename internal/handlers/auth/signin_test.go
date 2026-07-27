package auth

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	activeUsers map[string]*user.User
}

func (m *mockUserRepository) FindActiveByEmail(email string) (*user.User, error) {
	u, ok := m.activeUsers[email]
	if !ok {
		return nil, errors.New("user not found or inactive")
	}
	return u, nil
}

func (m *mockUserRepository) Create(u *user.User) error             { return nil }
func (m *mockUserRepository) Find(id uuid.UUID) (*user.User, error) { return nil, nil }
func (m *mockUserRepository) FindByEmail(e string) (*user.User, error) {
	return nil, nil
}
func (m *mockUserRepository) UpdatePassword(u *user.User) error { return nil }

func TestSignin(t *testing.T) {
	e := echo.New()
	e.Validator = &api.Validator{Validator: validator.New()}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	secret := "test-secret-key-32-chars-long-12345"
	config := &internal.Config{
		Secret: secret,
	}

	salt := "bUcdCORadqkbqHa1"
	plainPassword := "sholva123"

	// Construct combined password hash matching Auth.buildPassword logic
	combined := plainPassword + salt + secret
	if len(combined) > 71 {
		combined = combined[:71]
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(combined), bcrypt.DefaultCost)
	require.NoError(t, err)

	// Fixtures matching database/fixtures.sql
	activeUserTealc := &user.User{
		Id:       uuid.MustParse("aff84550-b21f-11ee-8ac0-5ab75f0c1cab"),
		Email:    "tealc.of.chulak@sgc.example.com",
		Name:     "Teal'c",
		Surname:  "Chulak",
		Salt:     salt,
		Password: string(hashedPassword),
		Role:     1,
	}

	mockRepo := &mockUserRepository{
		activeUsers: map[string]*user.User{
			"tealc.of.chulak@sgc.example.com": activeUserTealc,
			// jonas.quinn@sgc.example.com is deleted in fixtures -> omitted from activeUsers map
			// cameron.mitchell@sgc.example.com is unconfirmed in fixtures -> omitted from activeUsers map
		},
	}

	h := &Auth{
		echo:          e,
		configuration: config,
		auditlog:      logger,
		user: &userRepositories{
			repository: mockRepo,
		},
	}

	tests := []struct {
		name           string
		contentType    string
		payload        string
		expectedStatus int
		expectInBody   []string
	}{
		{
			name:           "Invalid JSON payload",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{invalid json`,
			expectedStatus: http.StatusUnauthorized,
			expectInBody:   []string{credentialError},
		},
		{
			name:           "Empty request body",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{}`,
			expectedStatus: http.StatusUnauthorized,
			expectInBody:   []string{credentialError},
		},
		{
			name:           "Missing password",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"tealc.of.chulak@sgc.example.com"}`,
			expectedStatus: http.StatusUnauthorized,
			expectInBody:   []string{credentialError},
		},
		{
			name:           "Invalid email format",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"not-an-email","password":"sholva123"}`,
			expectedStatus: http.StatusUnauthorized,
			expectInBody:   []string{credentialError},
		},
		{
			name:           "Non-existent user",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"unknown@sgc.example.com","password":"sholva123"}`,
			expectedStatus: http.StatusUnauthorized,
			expectInBody:   []string{credentialError},
		},
		{
			name:           "Deleted user fixture (jonas.quinn)",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"jonas.quinn@sgc.example.com","password":"sholva123"}`,
			expectedStatus: http.StatusUnauthorized,
			expectInBody:   []string{credentialError},
		},
		{
			name:           "Unconfirmed user fixture (cameron.mitchell)",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"cameron.mitchell@sgc.example.com","password":"sholva123"}`,
			expectedStatus: http.StatusUnauthorized,
			expectInBody:   []string{credentialError},
		},
		{
			name:           "Wrong password for active user fixture",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"tealc.of.chulak@sgc.example.com","password":"wrongpassword"}`,
			expectedStatus: http.StatusUnauthorized,
			expectInBody:   []string{"hashedpassword"},
		},
		{
			name:           "Valid credentials for active user fixture (tealc.of.chulak)",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"tealc.of.chulak@sgc.example.com","password":"sholva123"}`,
			expectedStatus: http.StatusOK,
			expectInBody:   []string{"token", "tealc.of.chulak@sgc.example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/auth/signin", bytes.NewBufferString(tt.payload))
			if tt.contentType != "" {
				req.Header.Set(echo.HeaderContentType, tt.contentType)
			}
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			err := h.signin(c)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			for _, expectedStr := range tt.expectInBody {
				assert.Contains(t, rec.Body.String(), expectedStr)
			}
		})
	}
}
