package handler

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/app"
	"github.com/jurgisjaska/binbogami/app/api"
	"github.com/jurgisjaska/binbogami/app/database"
	"github.com/labstack/echo/v4"
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

	h.echo.POST("/auth/signin", h.signin)
	h.echo.POST("/auth/signup", h.signup)
	h.echo.DELETE("/auth/signout", h.signout)

	return h
}

func (h *Auth) signin(c echo.Context) error {
	credentials := &api.Credentials{}
	if err := c.Bind(credentials); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect credentials"))
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(credentials); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect credentials", err.Error()))
	}

	// maybe there should be a service?
	// validate if password is correct

	user, err := h.repository.FindBy("email", credentials.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect credentials", err.Error()))
	}

	password := fmt.Sprintf("%s%s%s", credentials.Password, user.Salt, h.configuration.Salt)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password[:71]), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(hashedPassword, api.CreateRequest(c)))
}

func (h *Auth) signup(c echo.Context) error {
	// not sure if this is the best place for signup
	// but putting in the users handler feels more dirty

	user := &database.User{}
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect user"))
	}

	//

	return nil
}

func (h *Auth) signout(c echo.Context) error {
	return nil
}

// CreateAuth creates instance of the auth handler
// Differs from other handlers authentication require application configuration and sal
func CreateAuth(e *echo.Echo, d *sqlx.DB, c *app.Config) *Auth {
	return (&Auth{echo: e, database: d, configuration: c}).initialize()
}
