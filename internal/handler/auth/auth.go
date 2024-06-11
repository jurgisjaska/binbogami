package auth

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/database"
	"github.com/jurgisjaska/binbogami/internal/database/user"
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
		invitation     *database.InvitationRepository
		member         *database.MemberRepository
		organization   *database.OrganizationRepository
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
		configuration *user.ConfigurationRepository
		passwordReset *user.PasswordResetRepository
	}
)

func (h *Auth) initialize() *Auth {
	h.invitation = database.CreateInvitation(h.database)
	h.member = database.CreateMember(h.database)
	h.organization = database.CreateOrganization(h.database)

	h.userRepository = &userRepository{
		user:          user.CreateUser(h.database),
		configuration: user.CreateConfiguration(h.database),
		passwordReset: user.CreatePasswordReset(h.database),
	}

	h.echo.PUT("/auth/signin", h.signin)
	h.echo.POST("/auth/signup", h.signup)

	h.echo.POST("/auth/forgot-password", h.forgot)
	h.echo.POST("/auth/reset-password", h.reset)

	return h
}

// membership determine if a user is a member of any organization and return the organization information if true
// if member has multiple organization but no default he will be marked as member but will not have default organization
func (h *Auth) membership(u *user.User) (bool, *database.Organization) {
	m := false
	var organization *uuid.UUID

	members, err := h.member.ManyByUser(u)
	if err == nil && len(*members) != 0 {
		m = true

		if len(*members) > 1 {
			defaultConfiguration, err := h.userRepository.configuration.FindDefaultOrganization(u)
			if err == nil && defaultConfiguration != nil {
				defaultId, _ := uuid.Parse(defaultConfiguration.Value)
				organization = &defaultId
			}
		} else {
			organization = (*members)[0].OrganizationId
		}
	}

	o, _ := h.organization.FindById(organization)
	return m, o
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

// CreateAuth creates a new instance of the Auth handler and initializes it.
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
