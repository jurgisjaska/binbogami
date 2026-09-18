package finance

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models"
	"github.com/jurgisjaska/binbogami/internal/database/location"
	"github.com/labstack/echo/v5"
)

type Location struct {
	echo       *echo.Group
	database   *sqlx.DB
	repository location.LocationRepository
	auditlog   *slog.Logger
}

func (h *Location) initialize() *Location {
	h.repository = location.CreateLocation(h.database)

	h.echo.GET("/locations", h.index)
	h.echo.GET("/locations/:id", h.show)

	// stats
	// entries

	return h
}

func (h *Location) index(c *echo.Context) error {
	request := api.CreateRequest(c)

	var locations *location.Locations
	var t int
	var err error

	locations, t, err = h.repository.FindMany(request)
	if err != nil {
		h.auditlog.Warn("location index error: incorrect request", "error", err.Error())
		return c.JSON(http.StatusNotFound, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(locations, request, t))
}

func (h *Location) show(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.auditlog.Warn("location show error: incorrect location", "error", err.Error())
		return c.JSON(http.StatusBadRequest, api.Error("incorrect location"))
	}

	entity, err := h.repository.Find(id)
	if err != nil {
		h.auditlog.Warn("location show error: location not found", "error", err.Error(), "location_id", id)
		return c.JSON(http.StatusNotFound, api.Error("location not found"))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}

// @deprecated
func (h *Location) create(c *echo.Context) error {
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

// CreateLocation initializes the Location resource, sets up routes, and returns the created Location instance.
func CreateLocation(g *echo.Group, d *sqlx.DB, l *slog.Logger) *Location {
	return (&Location{echo: g, database: d, auditlog: l}).initialize()
}
