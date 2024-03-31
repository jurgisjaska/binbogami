package handler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/model"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/database"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
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

	h.echo.PUT("/auth", h.signin)
	h.echo.POST("/auth", h.signup)

	return h
}

// signin in creates new JWT token for the user if credentials are correct
func (h *Auth) signin(c echo.Context) error {
	request := &model.Signin{}
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(credentialError))
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors(credentialError, err.Error()))
	}

	u, err := h.user.By("email", request.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors(credentialError, err.Error()))
	}

	password := fmt.Sprintf("%s%s%s", request.Password, u.Salt, h.configuration.Secret)
	if len(password) > 71 {
		password = password[:71]
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return c.JSON(http.StatusInternalServerError, api.Error(err.Error()))
	}

	t, err := token.CreateToken(u, h.configuration.Secret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Error(err.Error()))
	}

	m, o := h.membership(u)
	response := model.SigninSuccess{Token: t, User: u, Member: m, Organization: o}

	return c.JSON(http.StatusOK, api.Success(response, api.CreateRequest(c)))
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
			if err != nil {
				defaultId, _ := uuid.Parse(defaultConfiguration.Value)
				organization = &defaultId
			}
		} else {
			organization = (*members)[0].OrganizationId
		}
	}

	o, _ := h.organization.ById(organization)
	return m, o
}

// signup validates signup form data and creates new user
// if the invitation UUID is present adds the new user to the organization
func (h *Auth) signup(c echo.Context) error {
	sm := &model.Signup{}
	if err := c.Bind(sm); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(signupError))
	}

	if err := c.Validate(sm); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors(signupError, err.Error()))
	}

	if sm.Password != sm.RepeatedPassword {
		return c.JSON(http.StatusBadRequest, api.Errors(signupError, fmt.Errorf("passwords does not match")))
	}

	existingUser, err := h.user.By("email", *sm.Email)
	if existingUser != nil {
		return c.JSON(http.StatusBadRequest, api.Error("email address already in use"))
	}

	u := &user.User{
		Email:   sm.Email,
		Name:    sm.Name,
		Surname: sm.Surname,
		Salt:    random.String(16),
	}

	u.Password, err = h.hashPassword(sm.Password, u.Salt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Errors(signupFailedError, err.Error()))
	}

	err = h.user.Create(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Errors(signupFailedError, err.Error()))
	}

	t, err := token.CreateToken(u, h.configuration.Secret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Error(err.Error()))
	}

	// @todo should invitation and user email match when using invitation link?

	member := &database.Member{}
	if sm.InvitationId != nil {
		invitation, err := h.invitation.Find(sm.InvitationId)
		if err == nil {
			member, err = h.member.Create(invitation.OrganizationId, u.Id, database.MemberRoleDefault, invitation.CreatedBy)
			if err == nil {
				_ = h.invitation.Delete(invitation)
			}
		}
	}

	isMember := false
	if member.Id != 0 {
		isMember = true
	}

	return c.JSON(
		http.StatusOK,
		api.Success(model.SignupSuccess{User: u, Token: t, IsMember: isMember}, api.CreateRequest(c)),
	)
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
