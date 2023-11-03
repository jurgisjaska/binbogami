package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type (
	Auth struct {
		echo     *echo.Echo
		database *sqlx.DB
	}
)

func (h *Auth) initialize() *Auth {
	h.echo.POST("/", h.signin)
	h.echo.DELETE("/", h.signout)

	return h
}

func (h *Auth) signin(c echo.Context) error {
	return nil
}

func (h *Auth) signout(c echo.Context) error {
	return nil
}

func CreateAuth(e *echo.Echo, d *sqlx.DB) *Auth {
	return (&Auth{echo: e, database: d}).initialize()
}
