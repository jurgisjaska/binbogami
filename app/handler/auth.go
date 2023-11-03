package handler

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/app/api"
	"github.com/jurgisjaska/binbogami/app/database"
	"github.com/labstack/echo/v4"
)

type (
	Auth struct {
		echo       *echo.Echo
		database   *sqlx.DB
		repository *database.UserRepository
	}
)

func (h *Auth) initialize() *Auth {
	h.repository = database.CreateUser(h.database)

	h.echo.POST("/", h.signin)
	h.echo.POST("/", h.signup)
	h.echo.DELETE("/", h.signout)

	return h
}

func (h *Auth) signin(c echo.Context) error {
	credentials := &struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if err := c.Bind(credentials); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect credentials"))
	}

	return nil
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

func CreateAuth(e *echo.Echo, d *sqlx.DB) *Auth {
	return (&Auth{echo: e, database: d}).initialize()
}
