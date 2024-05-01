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
	"github.com/labstack/gommon/log"
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

	h.echo.POST("/auth/forgot-password", h.forgot)
	h.echo.POST("/auth/reset-password", h.reset)

	return h
}

// signin in creates new JWT token for the user if credentials are correct
func (h *Auth) signin(c echo.Context) error {
	request := &model.Signin{}
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusUnauthorized, api.Error(credentialError))
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnauthorized, api.Errors(credentialError, err.Error()))
	}

	u, err := h.user.FindByColumn("email", request.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.Errors(credentialError, err.Error()))
	}

	password := fmt.Sprintf("%s%s%s", request.Password, u.Salt, h.configuration.Secret)
	if len(password) > 71 {
		password = password[:71]
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return c.JSON(http.StatusUnauthorized, api.Error(err.Error()))
	}

	t, err := token.CreateToken(u, h.configuration.Secret)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.Error(err.Error()))
	}

	m, o := h.membership(u)
	response := model.SigninResponse{Token: t, User: u, Member: m, Organization: o}

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

// signup validates signup form data and creates new user
// if the invitation UUID is present adds the new user to the organization
func (h *Auth) signup(c echo.Context) error {
	request := &model.SignupRequest{}
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(signupError))
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors(signupError, err.Error()))
	}

	if request.Password != request.RepeatedPassword {
		return c.JSON(http.StatusBadRequest, api.Errors(signupError, fmt.Errorf("passwords does not match")))
	}

	existingUser, err := h.user.FindByColumn("email", *request.Email)
	if existingUser != nil {
		return c.JSON(http.StatusBadRequest, api.Error("email address already in use"))
	}

	u := &user.User{
		Email:   request.Email,
		Name:    request.Name,
		Surname: request.Surname,
		Salt:    random.String(16),
	}

	u.Password, err = h.hashPassword(request.Password, u.Salt)
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
	organization := &database.Organization{}
	log.Infof("%+v", request)
	if request.InvitationId != nil {
		invitation, err := h.invitation.Find(request.InvitationId)
		log.Infof("%+v", invitation)
		log.Error(err)
		if err == nil {
			member, err = h.member.Create(invitation.OrganizationId, u.Id, database.MemberRoleDefault, invitation.CreatedBy)
			if err == nil {
				_ = h.invitation.Delete(invitation)
			}

			// error does not matter in this case
			// organization either is there or no
			// SQL not found can be ignored
			organization, _ = h.organization.FindById(invitation.OrganizationId)
		}
	}

	log.Infof("%+v", member)

	// membership status
	m := false
	if member.Id != 0 {
		m = true
	}

	// reset organization to nil
	// to keep consistency between sign in and sign up methods
	if organization.Id == nil {
		organization = nil
	}

	return c.JSON(
		http.StatusOK,
		api.Success(
			model.SignupResponse{User: u, Token: t, Member: m, Organization: organization},
			api.CreateRequest(c),
		),
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

func (h *Auth) forgot(c echo.Context) error {
	request := &model.ForgotPasswordRequest{}
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusUnauthorized, api.Error(credentialError))
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnauthorized, api.Errors(credentialError, err.Error()))
	}

	user, err := h.user.FindByColumn("email", request.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.Errors(credentialError, err.Error()))
	}

	// create reset token
	// send an email to user

	return nil
}

func (h *Auth) reset(c echo.Context) error {
	request := &model.ResetPasswordRequest{}
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusUnauthorized, api.Error(credentialError))
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnauthorized, api.Errors(credentialError, err.Error()))
	}

	// unique reset token is needed to reset password
	// it can be used once
	// expires after 24 hours

	return nil
}

// CreateAuth creates instance of the auth handler
func CreateAuth(e *echo.Echo, d *sqlx.DB, c *internal.Config) *Auth {
	return (&Auth{echo: e, database: d, configuration: c}).initialize()
}
