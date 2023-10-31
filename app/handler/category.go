package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/app/api"
	"github.com/jurgisjaska/binbogami/app/database"
	"github.com/labstack/echo/v4"
)

type (
	Category struct {
		echo       *echo.Echo
		database   *sqlx.DB
		repository *database.CategoryRepository
	}
)

func (h *Category) initialize() *Category {
	h.repository = database.CreateCategory(h.database)
	h.echo.GET("/categories/:id", h.one)
	h.echo.GET("/categories", h.many)

	return h
}

func (h *Category) one(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect category"))
	}

	category, err := h.repository.Find(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("category not found"))
	}

	return c.JSON(http.StatusOK, api.Success(category, api.CreateRequest(c)))
}

func (h *Category) many(c echo.Context) error {
	categories, err := h.repository.FindMany()
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("no categories found"))
	}

	return c.JSON(http.StatusOK, api.Success(categories, api.CreateRequest(c)))
}

// update
// create
// delete

func CreateCategory(e *echo.Echo, d *sqlx.DB) *Category {
	return (&Category{
		echo:     e,
		database: d,
	}).initialize()
}
