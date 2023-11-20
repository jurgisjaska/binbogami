package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/app"
	"github.com/jurgisjaska/binbogami/app/api"
	"github.com/jurgisjaska/binbogami/app/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"golang.org/x/crypto/bcrypt"
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

	// @todo probably should use POST, PUT, DELETE to single endpoint
	h.echo.POST("/auth/signin", h.signin)
	h.echo.POST("/auth/signup", h.signup)
	// h.echo.DELETE("/auth/signout", h.signout)

	return h
}

func (h *Auth) signin(c echo.Context) error {
	sm := &api.SigninModel{}
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

	password := fmt.Sprintf("%s%s%s", sm.Password, user.Salt, h.configuration.Salt)
	if len(password) > 71 {
		password = password[:71]
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.JSON(http.StatusInternalServerError, api.Error(err.Error()))
	}

	// generate token
	expire := jwt.NewNumericDate(time.Now().Add(time.Hour * 72))
	claim := &api.TokenClaims{
		Id:    *user.Id,
		Email: *user.Email,
		Name:  *user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expire,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return c.JSON(http.StatusOK, api.Success(token, api.CreateRequest(c)))
}

func (h *Auth) signup(c echo.Context) error {
	sm := &api.SignupModel{}
	if err := c.Bind(sm); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect signup information"))
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(sm); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect signup information", err.Error()))
	}

	user, err := h.repository.FindBy("email", sm.Email)
	if user != nil {
		return c.JSON(http.StatusBadRequest, api.Error("email address already in user"))
	}

	// @todo these pointer / no-pointer conversions feels incorrect for some reason.
	u := &database.User{}
	u.Email = &sm.Email
	u.Name = &sm.Name
	u.Surname = &sm.Surname
	u.Salt = random.String(16)
	u.Password, err = h.hashPassword(sm.Password, u.Salt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Errors("signup failed", err.Error()))
	}

	err = h.repository.Create(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Errors("signup failed", err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(u, api.CreateRequest(c)))
}

// hashPassword creates new password hash using bcrypt
func (h *Auth) hashPassword(password string, salt string) (string, error) {
	p := fmt.Sprintf("%s%s%s", password, salt, h.configuration.Salt)

	if len(p) > 71 {
		p = p[:71]
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// func (h *Auth) signout(c echo.Context) error {
// 	return nil
// }

// CreateAuth creates instance of the auth handler
// Differs from other handlers authentication require application configuration and sal
func CreateAuth(e *echo.Echo, d *sqlx.DB, c *app.Config) *Auth {
	return (&Auth{echo: e, database: d, configuration: c}).initialize()
}
