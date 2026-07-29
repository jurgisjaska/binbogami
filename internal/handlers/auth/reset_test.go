package auth

import (
	"bytes"
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
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/jurgisjaska/binbogami/internal/database/user/password"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReset(t *testing.T) {
	e := echo.New()
	e.Validator = &api.Validator{Validator: validator.New()}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	secret := "test-secret-key-32-chars-long-12345"
	config := &internal.Config{
		Secret: secret,
	}

	salt := "bUcdCORadqkbqHa1"

	// User fixtures matching database/fixtures.sql
	danielID := uuid.MustParse("2b63b228-b21c-11ee-9a7a-5ab75f0c1cab")
	tealcID := uuid.MustParse("aff84550-b21f-11ee-8ac0-5ab75f0c1cab")
	orphanUserID := uuid.MustParse("99999999-9999-9999-9999-999999999999")

	activeUserDaniel := &user.User{
		Id:       danielID,
		Email:    "daniel.jackson@sgc.example.com",
		Name:     "Daniel",
		Surname:  "Jackson",
		Salt:     salt,
		Password: "oldpasswordhash",
		Role:     1,
	}
	activeUserTealc := &user.User{
		Id:       tealcID,
		Email:    "tealc.of.chulak@sgc.example.com",
		Name:     "Teal'c",
		Surname:  "Chulak",
		Salt:     salt,
		Password: "oldpasswordhash",
		Role:     1,
	}

	// Password reset token fixtures matching database/fixtures.sql
	tokenDaniel := uuid.MustParse("27b8bb9a-dcfd-41cd-b60c-00f6cb7b89a1")
	tokenTealc := uuid.MustParse("aa7d1712-c6d0-488a-8b0d-0172f28b05d0")
	tokenOrphan := uuid.MustParse("88888888-8888-8888-8888-888888888888")
	tokenUnopened := uuid.MustParse("77777777-7777-7777-7777-777777777777")

	now := time.Now()
	resetDaniel := &password.Reset{
		Id:        tokenDaniel,
		UserId:    danielID,
		Ip:        "::1",
		UserAgent: "PostmanRuntime/7.51.0",
		CreatedAt: now,
		OpenedAt:  &now,
		ExpireAt:  now.Add(time.Hour * 2),
	}
	resetTealc := &password.Reset{
		Id:        tokenTealc,
		UserId:    tealcID,
		Ip:        "::1",
		UserAgent: "PostmanRuntime/7.51.0",
		CreatedAt: now,
		OpenedAt:  &now,
		ExpireAt:  now.Add(time.Hour * 2),
	}
	resetOrphan := &password.Reset{
		Id:        tokenOrphan,
		UserId:    orphanUserID,
		Ip:        "::1",
		UserAgent: "PostmanRuntime/7.51.0",
		CreatedAt: now,
		OpenedAt:  &now,
		ExpireAt:  now.Add(time.Hour * 2),
	}
	resetUnopened := &password.Reset{
		Id:        tokenUnopened,
		UserId:    danielID,
		Ip:        "::1",
		UserAgent: "PostmanRuntime/7.51.0",
		CreatedAt: now,
		OpenedAt:  nil,
		ExpireAt:  now.Add(time.Hour * 2),
	}

	tests := []struct {
		name                 string
		contentType          string
		payload              string
		failOnUpdatePassword bool
		failOnUpdateExpireAt bool
		expectedStatus       int
		expectInBody         []string
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
			name:           "Password too short",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"password":"short","repeatedPassword":"short","token":"27b8bb9a-dcfd-41cd-b60c-00f6cb7b89a1"}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectInBody:   []string{"error encountered during data validation"},
		},
		{
			name:           "Mismatched repeated password",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"password":"newpassword123","repeated_password":"differentpassword123","token":"27b8bb9a-dcfd-41cd-b60c-00f6cb7b89a1"}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectInBody:   []string{"error encountered during data validation"},
		},
		{
			name:           "Invalid token UUID format",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"password":"newpassword123","repeated_password":"newpassword123","token":"invalid-uuid-string"}`,
			expectedStatus: http.StatusBadRequest,
			expectInBody:   []string{"bad request"},
		},
		{
			name:           "Non-existent password reset token fixture",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"password":"newpassword123","repeated_password":"newpassword123","token":"11111111-1111-1111-1111-111111111111"}`,
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"password reset token not found"},
		},
		{
			name:           "Unopened password reset token fixture",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"password":"newpassword123","repeated_password":"newpassword123","token":"77777777-7777-7777-7777-777777777777"}`,
			expectedStatus: http.StatusUnauthorized,
			expectInBody:   []string{"password reset token not opened"},
		},
		{
			name:           "Token linked to non-existent user fixture",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"password":"newpassword123","repeated_password":"newpassword123","token":"88888888-8888-8888-8888-888888888888"}`,
			expectedStatus: http.StatusUnauthorized,
			expectInBody:   []string{"incorrect credentials"},
		},
		{
			name:                 "Database error updating user password",
			contentType:          echo.MIMEApplicationJSON,
			payload:              `{"password":"newpassword123","repeated_password":"newpassword123","token":"27b8bb9a-dcfd-41cd-b60c-00f6cb7b89a1"}`,
			failOnUpdatePassword: true,
			expectedStatus:       http.StatusInternalServerError,
			expectInBody:         []string{"internal server error"},
		},
		{
			name:                 "Database error updating password reset expire_at",
			contentType:          echo.MIMEApplicationJSON,
			payload:              `{"password":"newpassword123","repeated_password":"newpassword123","token":"27b8bb9a-dcfd-41cd-b60c-00f6cb7b89a1"}`,
			failOnUpdateExpireAt: true,
			expectedStatus:       http.StatusInternalServerError,
			expectInBody:         []string{"internal server error"},
		},
		{
			name:           "Successful password reset for active user fixture (daniel.jackson)",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"password":"newpassword123","repeated_password":"newpassword123","token":"27b8bb9a-dcfd-41cd-b60c-00f6cb7b89a1"}`,
			expectedStatus: http.StatusOK,
			expectInBody:   []string{"success", danielID.String(), "daniel.jackson@sgc.example.com"},
		},
		{
			name:           "Successful password reset for active user fixture (tealc.of.chulak)",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"password":"newpassword123","repeated_password":"newpassword123","token":"aa7d1712-c6d0-488a-8b0d-0172f28b05d0"}`,
			expectedStatus: http.StatusOK,
			expectInBody:   []string{"success", tealcID.String(), "tealc.of.chulak@sgc.example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &mockUserRepository{
				usersByID: map[uuid.UUID]*user.User{
					danielID: activeUserDaniel,
					tealcID:  activeUserTealc,
				},
				failOnUpdatePassword: tt.failOnUpdatePassword,
			}

			mockResetRepo := &mockPasswordResetRepository{
				resetsByID: map[uuid.UUID]*password.Reset{
					tokenDaniel:   resetDaniel,
					tokenTealc:    resetTealc,
					tokenOrphan:   resetOrphan,
					tokenUnopened: resetUnopened,
				},
				failOnUpdate: tt.failOnUpdateExpireAt,
			}

			h := &Auth{
				echo:          e,
				configuration: config,
				auditlog:      logger,
				user: &userRepositories{
					repository:    mockUserRepo,
					passwordReset: mockResetRepo,
				},
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/reset-password", bytes.NewBufferString(tt.payload))
			if tt.contentType != "" {
				req.Header.Set(echo.HeaderContentType, tt.contentType)
			}
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			err := h.reset(c)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			for _, expectedStr := range tt.expectInBody {
				assert.Contains(t, rec.Body.String(), expectedStr)
			}
		})
	}
}

func TestOpen(t *testing.T) {
	e := echo.New()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	secret := "test-secret-key-32-chars-long-12345"
	config := &internal.Config{
		Secret: secret,
	}

	tokenDaniel := uuid.MustParse("27b8bb9a-dcfd-41cd-b60c-00f6cb7b89a1")

	resetDaniel := &password.Reset{
		Id:        tokenDaniel,
		UserId:    uuid.MustParse("2b63b228-b21c-11ee-9a7a-5ab75f0c1cab"),
		Ip:        "::1",
		UserAgent: "PostmanRuntime/7.51.0",
		CreatedAt: time.Now(),
		ExpireAt:  time.Now().Add(time.Hour * 2),
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
			name:           "Invalid password reset token UUID format",
			paramID:        "invalid-uuid-string",
			expectedStatus: http.StatusBadRequest,
			expectInBody:   []string{"incorrect password reset token"},
		},
		{
			name:           "Password reset token not found",
			paramID:        "11111111-1111-1111-1111-111111111111",
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"password reset token not found"},
		},
		{
			name:           "Database error updating password reset token",
			paramID:        "27b8bb9a-dcfd-41cd-b60c-00f6cb7b89a1",
			failOnUpdate:   true,
			expectedStatus: http.StatusInternalServerError,
			expectInBody:   []string{"failed to update password reset token"},
		},
		{
			name:           "Successful password reset token open",
			paramID:        "27b8bb9a-dcfd-41cd-b60c-00f6cb7b89a1",
			expectedStatus: http.StatusOK,
			expectInBody:   []string{"success", tokenDaniel.String()},
			checkOpenedAt:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetCopy := *resetDaniel
			mockResetRepo := &mockPasswordResetRepository{
				resetsByID: map[uuid.UUID]*password.Reset{
					tokenDaniel: &resetCopy,
				},
				failOnUpdate: tt.failOnUpdate,
			}

			h := &Auth{
				echo:          e,
				configuration: config,
				auditlog:      logger,
				user: &userRepositories{
					passwordReset: mockResetRepo,
				},
			}

			req := httptest.NewRequest(http.MethodGet, "/auth/reset-password/"+tt.paramID, nil)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath("/auth/reset-password/:id")
			c.SetPathValues(echo.PathValues{echo.PathValue{Name: "id", Value: tt.paramID}})

			err := h.open(c)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			for _, expectedStr := range tt.expectInBody {
				assert.Contains(t, rec.Body.String(), expectedStr)
			}

			if tt.checkOpenedAt {
				assert.NotNil(t, resetCopy.OpenedAt)
			}
		})
	}
}
