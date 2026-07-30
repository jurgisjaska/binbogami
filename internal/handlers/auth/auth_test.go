package auth

import (
	"io"
	"log/slog"
	"net/http"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

func TestBuildPassword(t *testing.T) {
	tests := []struct {
		name     string
		secret   string
		password string
		salt     string
		expected string
	}{
		{
			name:     "Combined length less than 71 characters",
			secret:   "secret123",
			password: "pass",
			salt:     "salt",
			expected: "passsaltsecret123",
		},
		{
			name:     "Combined length exactly 71 characters",
			secret:   "1234567890123456789012345678901234567890123456789012345678901234567",
			password: "pa",
			salt:     "sa",
			expected: "pasa1234567890123456789012345678901234567890123456789012345678901234567",
		},
		{
			name:     "Combined length greater than 71 characters is truncated",
			secret:   "secret-key-that-is-very-long-and-exceeds-the-maximum-length-allowed-by-bcrypt-processing-limits",
			password: "password123",
			salt:     "salt456",
			expected: "password123salt456secret-key-that-is-very-long-and-exceeds-the-maximum-",
		},
		{
			name:     "Empty password and salt",
			secret:   "mysecret",
			password: "",
			salt:     "",
			expected: "mysecret",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Auth{
				configuration: &internal.Config{
					Secret: tt.secret,
				},
			}
			result := h.buildPassword(tt.password, tt.salt)
			assert.Equal(t, tt.expected, result)
			assert.LessOrEqual(t, len(result), 71)
		})
	}
}

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		secret   string
		password string
		salt     string
	}{
		{
			name:     "Standard password and salt",
			secret:   "testsecret32byteslongkeyforauth!",
			password: "mysecretpassword123",
			salt:     "randomsalt123",
		},
		{
			name:     "Long password and salt resulting in truncation",
			secret:   "super-long-secret-key-exceeding-the-seventy-one-character-limitations-of-bcrypt",
			password: "anotherpassword456",
			salt:     "anothersalt789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Auth{
				configuration: &internal.Config{
					Secret: tt.secret,
				},
			}
			hashedPassword, err := h.hashPassword(tt.password, tt.salt)
			require.NoError(t, err)
			assert.NotEmpty(t, hashedPassword)

			built := h.buildPassword(tt.password, tt.salt)
			err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(built))
			assert.NoError(t, err)
		})
	}
}

func TestInitialize(t *testing.T) {
	e := echo.New()
	db := &sqlx.DB{}
	cfg := &internal.Config{Secret: "test-secret"}

	h := &Auth{
		echo:          e,
		database:      db,
		configuration: cfg,
	}

	res := h.initialize()
	assert.Same(t, h, res)
	assert.NotNil(t, h.invitation)
	assert.NotNil(t, h.user)
	assert.NotNil(t, h.user.repository)
	assert.NotNil(t, h.user.configuration)
	assert.NotNil(t, h.user.passwordReset)

	routes := e.Router().Routes()
	expectedRoutes := []struct {
		method string
		path   string
	}{
		{http.MethodPut, "/auth/signin"},
		{http.MethodPost, "/auth/signup"},
		{http.MethodPost, "/auth/forgot-password"},
		{http.MethodGet, "/auth/reset-password/:id"},
		{http.MethodPost, "/auth/reset-password"},
	}

	for _, er := range expectedRoutes {
		found := false
		for _, r := range routes {
			if r.Method == er.method && r.Path == er.path {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected route %s %s not registered", er.method, er.path)
	}
}

func TestCreateAuth(t *testing.T) {
	e := echo.New()
	db := &sqlx.DB{}
	cfg := &internal.Config{Secret: "test-secret"}
	dialer := &gomail.Dialer{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	auth := CreateAuth(e, db, cfg, dialer, logger)

	require.NotNil(t, auth)
	assert.Same(t, e, auth.echo)
	assert.Same(t, db, auth.database)
	assert.Same(t, cfg, auth.configuration)
	assert.Same(t, logger, auth.auditlog)
	assert.NotNil(t, auth.mailer)
	assert.NotNil(t, auth.mailer.resetPassword)
	assert.NotNil(t, auth.invitation)
	assert.NotNil(t, auth.user)

	routes := e.Router().Routes()
	assert.Len(t, routes, 5)
}
