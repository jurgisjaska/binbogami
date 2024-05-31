package auth

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/database"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	signupError       string = "incorrect signup information"
	credentialError   string = "incorrect credentials"
	signupFailedError string = "signup failed"
)

type (
	Auth struct {
		echo              *echo.Echo
		database          *sqlx.DB
		user              *user.Repository
		userConfiguration *user.ConfigurationRepository
		invitation        *database.InvitationRepository
		member            *database.MemberRepository
		organization      *database.OrganizationRepository
		configuration     *internal.Config
	}
)

func (h *Auth) initialize() *Auth {
	h.user = user.CreateUser(h.database)
	h.invitation = database.CreateInvitation(h.database)
	h.member = database.CreateMember(h.database)
	h.userConfiguration = user.CreateConfiguration(h.database)
	h.organization = database.CreateOrganization(h.database)

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
			defaultConfiguration, err := h.userConfiguration.DefaultOrganization(u)
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

// CreateAuth creates instance of the auth handler
func CreateAuth(e *echo.Echo, d *sqlx.DB, c *internal.Config) *Auth {
	return (&Auth{echo: e, database: d, configuration: c}).initialize()
}
