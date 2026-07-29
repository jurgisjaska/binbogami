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
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/jurgisjaska/binbogami/internal/database/user/password"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPasswordResetRepository struct {
	resets               map[uuid.UUID]password.Resets
	resetsByID           map[uuid.UUID]*password.Reset
	failOnCreate         bool
	failOnUpdateExpireAt bool
	failOnUpdate         bool
}

func (m *mockPasswordResetRepository) Create(pr *password.Reset) error {
	if m.failOnCreate {
		return errors.New("database error inserting password reset")
	}
	if m.resets == nil {
		m.resets = make(map[uuid.UUID]password.Resets)
	}
	m.resets[pr.UserId] = append(m.resets[pr.UserId], *pr)
	return nil
}

func (m *mockPasswordResetRepository) Find(id uuid.UUID) (*password.Reset, error) {
	if m.resetsByID != nil {
		if pr, ok := m.resetsByID[id]; ok {
			return pr, nil
		}
	}
	return nil, errors.New("password reset token not found")
}

func (m *mockPasswordResetRepository) FindManyByUser(u *user.User, limit int) (*password.Resets, error) {
	res, ok := m.resets[u.Id]
	if !ok {
		empty := password.Resets{}
		return &empty, nil
	}
	return &res, nil
}

func (m *mockPasswordResetRepository) UpdateExpireAt(u *user.User) error {
	if m.failOnUpdateExpireAt {
		return errors.New("database error updating expire_at")
	}
	return nil
}

func (m *mockPasswordResetRepository) Update(pr *password.Reset) error {
	if m.failOnUpdate {
		return errors.New("database error updating password reset")
	}
	if m.resetsByID != nil {
		m.resetsByID[pr.Id] = pr
	}
	return nil
}

type mockResetPasswordMailer struct {
	failOnSend bool
}

func (m *mockResetPasswordMailer) Send(u *user.User, pr *password.Reset) error {
	if m.failOnSend {
		return errors.New("smtp server connection failed")
	}
	return nil
}

func TestForgot(t *testing.T) {
	e := echo.New()
	e.Validator = &api.Validator{Validator: validator.New()}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	secret := "test-secret-key-32-chars-long-12345"
	config := &internal.Config{
		Secret: secret,
	}

	// User fixtures matching database/fixtures.sql
	tealcID := uuid.MustParse("aff84550-b21f-11ee-8ac0-5ab75f0c1cab")
	jackID := uuid.MustParse("05e7257a-b21c-11ee-9a7a-5ab75f0c1cab")
	samanthaID := uuid.MustParse("1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab")
	danielID := uuid.MustParse("2b63b228-b21c-11ee-9a7a-5ab75f0c1cab")

	activeUserTealc := &user.User{
		Id:      tealcID,
		Email:   "tealc.of.chulak@sgc.example.com",
		Name:    "Teal'c",
		Surname: "Chulak",
		Role:    1,
	}
	activeUserJack := &user.User{
		Id:      jackID,
		Email:   "jack.oneil@sgc.example.com",
		Name:    "Jack",
		Surname: "O'Neil",
		Role:    6,
	}
	activeUserSam := &user.User{
		Id:      samanthaID,
		Email:   "samantha.carter@sgc.example.com",
		Name:    "Samantha",
		Surname: "Carter",
		Role:    4,
	}
	activeUserDaniel := &user.User{
		Id:      danielID,
		Email:   "daniel.jackson@sgc.example.com",
		Name:    "Daniel",
		Surname: "Jackson",
		Role:    1,
	}

	mockUserRepo := &mockUserRepository{
		activeUsers: map[string]*user.User{
			"tealc.of.chulak@sgc.example.com": activeUserTealc,
			"jack.oneil@sgc.example.com":      activeUserJack,
			"samantha.carter@sgc.example.com": activeUserSam,
			"daniel.jackson@sgc.example.com":  activeUserDaniel,
			// jonas.quinn@sgc.example.com is deleted in fixtures -> omitted from activeUsers map
			// cameron.mitchell@sgc.example.com is unconfirmed in fixtures -> omitted from activeUsers map
		},
	}

	// 10 existing password reset fixtures for Daniel Jackson (matching database/fixtures.sql)
	danielResets := password.Resets{}
	for i := 0; i < 10; i++ {
		danielResets = append(danielResets, password.Reset{
			Id:        uuid.New(),
			UserId:    danielID,
			Ip:        "::1",
			UserAgent: "PostmanRuntime/7.51.0",
			CreatedAt: time.Now(),
			ExpireAt:  time.Now().Add(time.Hour * 2),
		})
	}

	tests := []struct {
		name           string
		contentType    string
		payload        string
		failCreateDB   bool
		failMailSend   bool
		expectedStatus int
		expectInBody   []string
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
			payload:        `{"email":"not-an-email"}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectInBody:   []string{"error encountered during data validation"},
		},
		{
			name:           "Non-existent user",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"unknown@sgc.example.com"}`,
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"user not found"},
		},
		{
			name:           "Deleted user fixture (jonas.quinn)",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"jonas.quinn@sgc.example.com"}`,
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"user not found"},
		},
		{
			name:           "Unconfirmed user fixture (cameron.mitchell)",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"cameron.mitchell@sgc.example.com"}`,
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"user not found"},
		},
		{
			name:           "Too many reset requests fixture (daniel.jackson)",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"daniel.jackson@sgc.example.com"}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectInBody:   []string{"too many password resets"},
		},
		{
			name:           "Database error on reset creation",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"jack.oneil@sgc.example.com"}`,
			failCreateDB:   true,
			expectedStatus: http.StatusInternalServerError,
			expectInBody:   []string{"database error inserting password reset"},
		},
		{
			name:           "Mailer failure on sending reset password email",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"samantha.carter@sgc.example.com"}`,
			failMailSend:   true,
			expectedStatus: http.StatusInternalServerError,
			expectInBody:   []string{"smtp server connection failed"},
		},
		{
			name:           "Successful password reset request for active user fixture (tealc.of.chulak)",
			contentType:    echo.MIMEApplicationJSON,
			payload:        `{"email":"tealc.of.chulak@sgc.example.com"}`,
			expectedStatus: http.StatusOK,
			expectInBody:   []string{"success", tealcID.String()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockResetRepo := &mockPasswordResetRepository{
				resets: map[uuid.UUID]password.Resets{
					danielID: danielResets,
				},
				failOnCreate: tt.failCreateDB,
			}
			mockMailer := &mockResetPasswordMailer{
				failOnSend: tt.failMailSend,
			}

			h := &Auth{
				echo:          e,
				configuration: config,
				auditlog:      logger,
				user: &userRepositories{
					repository:    mockUserRepo,
					passwordReset: mockResetRepo,
				},
				mailer: &mailer{
					resetPassword: mockMailer,
				},
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/forgot-password", bytes.NewBufferString(tt.payload))
			if tt.contentType != "" {
				req.Header.Set(echo.HeaderContentType, tt.contentType)
			}
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			err := h.forgot(c)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			for _, expectedStr := range tt.expectInBody {
				assert.Contains(t, rec.Body.String(), expectedStr)
			}
		})
	}
}
