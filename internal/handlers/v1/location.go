package v1

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models"
	"github.com/jurgisjaska/binbogami/internal/database/book"
	"github.com/jurgisjaska/binbogami/internal/database/location"
	"github.com/labstack/echo/v4"
)

type Location struct {
	echo       *echo.Group
	database   *sqlx.DB
	repository *location.LocationRepository
	book       *book.Repository
}

func (h *Location) initialize() *Location {
	h.repository = location.CreateLocation(h.database)
	h.book = book.CreateBook(h.database)

	h.echo.POST("/locations", h.create)
	h.echo.GET("/books/:id/locations", h.byBook)

	return h
}

func (h *Location) one(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect location"))
	}

	location, err := h.repository.Find(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("location not found"))
	}

	return c.JSON(http.StatusOK, api.Success(location, api.CreateRequest(c)))
}

func (h *Location) create(c echo.Context) error {
	location := &models.Location{}
	if err := c.Bind(location); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect location data"))
	}

	if err := c.Validate(location); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect location data", err.Error()))
	}

	// location.CreatedBy = member.UserId
	entity, err := h.repository.Create(location)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}

func (h *Location) byBook(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect book"))
	}

	book, err := h.book.Find(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("book not found"))
	}

	locations, err := h.repository.ManyByBook(book)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("no locations found"))
	}

	return c.JSON(http.StatusOK, api.Success(locations, api.CreateRequest(c)))
}

func CreateLocation(g *echo.Group, d *sqlx.DB) *Location {
	return (&Location{echo: g, database: d}).initialize()
}
