package user

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/labstack/echo/v4"
)

type User struct {
	echo       *echo.Group
	database   *sqlx.DB
	repository *user.Repository
}

func (h *User) initialize() *User {
	h.repository = user.CreateUser(h.database)

	h.echo.GET("/users/:id", h.one)
	h.echo.GET("/users", h.many)

	return h
}

func (h *User) one(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect user"))
	}

	user, err := h.repository.By("id", id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("user not found"))
	}

	return c.JSON(http.StatusOK, api.Success(user, api.CreateRequest(c)))
}

func (h *User) many(c echo.Context) error {
	filter := c.QueryParam("filter")
	users, err := h.repository.FindMany(filter)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("no users found"))
	}

	return c.JSON(http.StatusOK, api.Success(users, api.CreateRequest(c)))
}

func CreateUser(g *echo.Group, d *sqlx.DB) *User {
	return (&User{echo: g, database: d}).initialize()
}
