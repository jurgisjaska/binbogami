package auth

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/jurgisjaska/binbogami/internal/database/user/configuration"
	"github.com/jurgisjaska/binbogami/internal/database/user/invitation"
	"github.com/jurgisjaska/binbogami/internal/database/user/password"
	"github.com/jurgisjaska/binbogami/internal/service/mail"
	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

const (
	credentialError string = "incorrect credentials"
	validationError string = "error encountered during data validation"
	requestError    string = "bad request"
	internalError   string = "internal server error"
)

type (
	Auth struct {
		echo          *echo.Echo
		database      *sqlx.DB
		invitation    *invitation.InvitationRepository
		configuration *internal.Config
		mailer        *mailer
		user          *userRepositories
		auditlog      *slog.Logger
	}

	// @todo go level up on a tree if there will not be any other mailers
	resetPasswordMailer interface {
		Send(u *user.User, pr *password.Reset) error
	}

	// @todo go level up on a tree if there will not be any other mailers
	mailer struct {
		resetPassword resetPasswordMailer
	}

	userRepositories struct {
		repository    user.UserRepository
		configuration *configuration.ConfigurationRepository
		passwordReset password.PasswordResetRepository
	}
)

func (h *Auth) initialize() *Auth {
	h.invitation = invitation.CreateInvitation(h.database)
	h.user = &userRepositories{
		repository:    user.CreateUser(h.database),
		configuration: configuration.CreateConfiguration(h.database),
		passwordReset: password.CreatePasswordReset(h.database),
	}

	h.echo.PUT("/auth/signin", h.signin)
	h.echo.POST("/auth/signup", h.signup)

	h.echo.POST("/auth/forgot-password", h.forgot)
	h.echo.GET("/auth/reset-password/:id", h.open)
	h.echo.POST("/auth/reset-password", h.reset)

	return h
}

// hashPassword creates a new password hash using bcrypt.
func (h *Auth) hashPassword(password string, salt string) (string, error) {
	p := h.buildPassword(password, salt)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (h *Auth) buildPassword(password string, salt string) string {
	p := fmt.Sprintf("%s%s%s", password, salt, h.configuration.Secret)

	if len(p) > 71 {
		p = p[:71]
	}

	return p
}

// CreateAuth creates a new instance of the Auth handlers and initializes it.
func CreateAuth(e *echo.Echo, d *sqlx.DB, c *internal.Config, md *gomail.Dialer, auditlog *slog.Logger) *Auth {
	return (&Auth{
		echo:          e,
		database:      d,
		configuration: c,
		mailer: &mailer{
			resetPassword: mail.CreateResetPassword(md, c),
		},
		auditlog: auditlog,
	}).initialize()
}
