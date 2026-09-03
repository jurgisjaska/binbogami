package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserRepository struct {
	user *user.User
	err  error
}

func (m *mockUserRepository) Find(id uuid.UUID) (*user.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.user, nil
}

func (m *mockUserRepository) FindActive(id uuid.UUID) (*user.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.user != nil && (m.user.DeletedAt != nil || m.user.ConfirmedAt == nil) {
		return nil, errors.New("user not active")
	}
	return m.user, nil
}

func (m *mockUserRepository) FindActiveByEmail(e string) (*user.User, error)     { return nil, nil }
func (m *mockUserRepository) FindNotDeletedByEmail(e string) (*user.User, error) { return nil, nil }
func (m *mockUserRepository) FindByEmail(e string) (*user.User, error)           { return nil, nil }
func (m *mockUserRepository) FindMany(filter string) (*user.Users, error)        { return nil, nil }
func (m *mockUserRepository) Create(u *user.User) error                          { return nil }
func (m *mockUserRepository) UpdatePassword(u *user.User) error                  { return nil }

func TestUserActive(t *testing.T) {
	now := time.Now()
	userId := uuid.New()

	tests := []struct {
		name           string
		tokenSetup     func(c *echo.Context)
		mockUser       *user.User
		mockErr        error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "valid active user",
			tokenSetup: func(c *echo.Context) {
				claims := &token.Claims{Id: &userId}
				tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				c.Set("user", tkn)
			},
			mockUser: &user.User{
				Id:          userId,
				Email:       "test@example.com",
				ConfirmedAt: &now,
				DeletedAt:   nil,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "ok",
		},
		{
			name: "unconfirmed user",
			tokenSetup: func(c *echo.Context) {
				claims := &token.Claims{Id: &userId}
				tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				c.Set("user", tkn)
			},
			mockUser: &user.User{
				Id:          userId,
				Email:       "unconfirmed@example.com",
				ConfirmedAt: nil,
				DeletedAt:   nil,
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `"user not confirmed"`,
		},
		{
			name: "deleted user",
			tokenSetup: func(c *echo.Context) {
				claims := &token.Claims{Id: &userId}
				tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				c.Set("user", tkn)
			},
			mockUser: &user.User{
				Id:          userId,
				Email:       "deleted@example.com",
				ConfirmedAt: &now,
				DeletedAt:   &now,
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `"user deleted"`,
		},
		{
			name: "user not found in repository",
			tokenSetup: func(c *echo.Context) {
				claims := &token.Claims{Id: &userId}
				tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				c.Set("user", tkn)
			},
			mockErr:        errors.New("user not found"),
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `"user not found"`,
		},
		{
			name: "missing user token in context",
			tokenSetup: func(c *echo.Context) {
				// No "user" key set
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `"unauthorized"`,
		},
		{
			name: "nil user ID in token claims",
			tokenSetup: func(c *echo.Context) {
				claims := &token.Claims{Id: nil}
				tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				c.Set("user", tkn)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `"unauthorized"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tt.tokenSetup(c)

			repo := &mockUserRepository{
				user: tt.mockUser,
				err:  tt.mockErr,
			}

			handler := UserActive(repo)(func(c *echo.Context) error {
				return c.String(http.StatusOK, "ok")
			})

			err := handler(c)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.expectedBody)
		})
	}
}

func TestCreateUserActive(t *testing.T) {
	repo := &mockUserRepository{}
	mw := CreateUserActive(repo)
	assert.NotNil(t, mw)
}
