package auth

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/jurgisjaska/binbogami/internal/database/user/configuration"
	"github.com/jurgisjaska/binbogami/internal/database/user/invitation"
	"github.com/jurgisjaska/binbogami/internal/database/user/password"
	"github.com/jurgisjaska/binbogami/internal/service/mail"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

const (
	credentialError     string = "incorrect credentials"
	validationError     string = "error encountered during data validation"
	requestError        string = "bad request"
	passwordsMatchError string = "passwords do not match"
	internalError       string = "internal server error"
)

type (
	Auth struct {
		echo           *echo.Echo
		database       *sqlx.DB
		invitation     *invitation.InvitationRepository
		configuration  *internal.Config
		mailer         *mailer
		userRepository *userRepository
	}

	// @todo go level up on a tree if there will not be any other mailers
	mailer struct {
		resetPassword *mail.ResetPassword
	}

	userRepository struct {
		user          *user.Repository
		configuration *configuration.ConfigurationRepository
		passwordReset *password.PasswordResetRepository
	}
)

func (h *Auth) initialize() *Auth {
	h.invitation = invitation.CreateInvitation(h.database)
	h.userRepository = &userRepository{
		user:          user.CreateUser(h.database),
		configuration: configuration.CreateConfiguration(h.database),
		passwordReset: password.CreatePasswordReset(h.database),
	}

	h.echo.PUT("/auth/signin", h.signin)
	h.echo.POST("/auth/signup", h.signup)

	h.echo.POST("/auth/forgot-password", h.forgot)
	h.echo.POST("/auth/reset-password", h.reset)

	return h
}

// hashPassword creates new password hash using bcrypt.
func (h *Auth) hashPassword(password string, salt string) (string, error) {
	p := fmt.Sprintf("%s%s%s", password, salt, h.configuration.Secret)

	if len(p) > 71 {
		p = p[:71]
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// CreateAuth creates a new instance of the Auth handlers and initializes it.
func CreateAuth(e *echo.Echo, d *sqlx.DB, c *internal.Config, md *gomail.Dialer) *Auth {
	return (&Auth{
		echo:          e,
		database:      d,
		configuration: c,
		mailer: &mailer{
			resetPassword: mail.CreateResetPassword(md, c),
		},
	}).initialize()
}
