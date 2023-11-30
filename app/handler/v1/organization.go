package v1

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/app/api"
	"github.com/jurgisjaska/binbogami/app/database"
	"github.com/labstack/echo/v4"
)

type (
	Organization struct {
		echo       *echo.Echo
		database   *sqlx.DB
		repository *database.OrganizationRepository
	}
)

func (h *Organization) initialize() *Organization {
	h.repository = database.CreateOrganization(h.database)
	h.echo.GET("/organizations/:id", h.one)
	// h.echo.GET("/organizations", h.many)
	// h.echo.POST("/organizations", h.create)
	// h.echo.PUT("/organizations/:id", h.update)
	// h.echo.DELETE("/organizations/:id", h.delete)

	return h
}

func (h *Organization) one(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect organization"))
	}

	organization, err := h.repository.Find(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("organization not found"))
	}

	return c.JSON(http.StatusOK, api.Success(organization, api.CreateRequest(c)))
}

func CreateOrganization(e *echo.Echo, d *sqlx.DB) *Organization {
	return (&Organization{echo: e, database: d}).initialize()
}
