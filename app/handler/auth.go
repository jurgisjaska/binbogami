package handler

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/app"
	"github.com/jurgisjaska/binbogami/app/api"
	"github.com/jurgisjaska/binbogami/app/api/token"
	"github.com/jurgisjaska/binbogami/app/database"
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
		repository    *database.UserRepository
		configuration *app.Config
	}
)

func (h *Auth) initialize() *Auth {
	h.repository = database.CreateUser(h.database)

	h.echo.PUT("/auth", h.signin)
	h.echo.POST("/auth", h.signup)

	return h
}

// signin in creates new JWT token for the user if credentials are correct
func (h *Auth) signin(c echo.Context) error {
	sm := &api.Signin{}
	if err := c.Bind(sm); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect credentials"))
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(sm); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect credentials", err.Error()))
	}

	user, err := h.repository.FindBy("email", sm.Email)
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

	return c.JSON(http.StatusOK, api.Success(api.SigninSuccess{t}, api.CreateRequest(c)))
}

// signup validates signup form data and creates new user
func (h *Auth) signup(c echo.Context) error {
	sm := &api.Signup{}
	if err := c.Bind(sm); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(signupError))
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(sm); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors(signupError, err.Error()))
	}

	// verify that passwords match
	if sm.Password != sm.RepeatedPassword {
		return c.JSON(http.StatusBadRequest, api.Errors(signupError, fmt.Errorf("passwords does not match")))
	}

	existingUser, err := h.repository.FindBy("email", *sm.Email)
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

	err = h.repository.Create(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Errors("signup failed", err.Error()))
	}

	t, err := token.CreateToken(u, h.configuration.Secret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(api.SignupSuccess{u, t}, api.CreateRequest(c)))
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
func CreateAuth(e *echo.Echo, d *sqlx.DB, c *app.Config) *Auth {
	return (&Auth{echo: e, database: d, configuration: c}).initialize()
}
