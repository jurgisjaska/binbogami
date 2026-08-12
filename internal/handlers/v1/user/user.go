package user

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/labstack/echo/v5"
)

type User struct {
	echo       *echo.Group
	database   *sqlx.DB
	repository *user.Repository
	auditlog   *slog.Logger
}

func (h *User) initialize() *User {
	h.repository = user.CreateUser(h.database)

	h.echo.GET("/users", h.index)
	h.echo.GET("/users/:id", h.show)

	return h
}

func (h *User) show(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect user"))
	}

	user, err := h.repository.Find(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("user not found"))
	}

	return c.JSON(http.StatusOK, api.Success(user, api.CreateRequest(c)))
}

func (h *User) index(c *echo.Context) error {
	filter := c.QueryParam("filter")
	users, err := h.repository.FindMany(filter)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("no users found"))
	}

	return c.JSON(http.StatusOK, api.Success(users, api.CreateRequest(c)))
}

func CreateUser(g *echo.Group, d *sqlx.DB, l *slog.Logger) *User {
	return (&User{echo: g, database: d, auditlog: l}).initialize()
}
