package handler

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/model"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"golang.org/x/crypto/bcrypt"
)

const (
	signupError string = "incorrect signup information"
)

type (
	Auth struct {
		echo          *echo.Echo
		database      *sqlx.DB
		user          *database.UserRepository
		invitation    *database.InvitationRepository
		member        *database.MemberRepository
		configuration *internal.Config
	}
)

func (h *Auth) initialize() *Auth {
	h.user = database.CreateUser(h.database)
	h.invitation = database.CreateInvitation(h.database)
	h.member = database.CreateMember(h.database)

	h.echo.PUT("/auth", h.signin)
	h.echo.POST("/auth", h.signup)

	return h
}

// signin in creates new JWT token for the user if credentials are correct
func (h *Auth) signin(c echo.Context) error {
	sm := &model.Signin{}
	if err := c.Bind(sm); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect credentials"))
	}

	if err := c.Validate(sm); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect credentials", err.Error()))
	}

	user, err := h.user.FindBy("email", sm.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect credentials", err.Error()))
	}

	password := fmt.Sprintf("%s%s%s", sm.Password, user.Salt, h.configuration.Secret)
	if len(password) > 71 {
		password = password[:71]
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.JSON(http.StatusInternalServerError, api.Error(err.Error()))
	}

	t, err := token.CreateToken(user, h.configuration.Secret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(model.SigninSuccess{t}, api.CreateRequest(c)))
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

	existingUser, err := h.user.FindBy("email", *sm.Email)
	if existingUser != nil {
		return c.JSON(http.StatusBadRequest, api.Error("email address already in use"))
	}

	u := &database.User{
		Email:   sm.Email,
		Name:    sm.Name,
		Surname: sm.Surname,
		Salt:    random.String(16),
	}

	u.Password, err = h.hashPassword(sm.Password, u.Salt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Errors("signup failed", err.Error()))
	}

	err = h.user.Create(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Errors("signup failed", err.Error()))
	}

	t, err := token.CreateToken(u, h.configuration.Secret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Error(err.Error()))
	}

	if sm.InvitationId != nil {
		invitation, err := h.invitation.Find(sm.InvitationId)
		if err == nil {
			_, err := h.member.Create(invitation.OrganizationId, u.Id, database.MemberRoleDefault, invitation.CreatedBy)
			if err == nil {
				_ = h.invitation.Delete(invitation)
			}
		}
	}

	return c.JSON(http.StatusOK, api.Success(model.SignupSuccess{u, t}, api.CreateRequest(c)))
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
