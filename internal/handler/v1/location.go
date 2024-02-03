package v1

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/model"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/database"
	"github.com/labstack/echo/v4"
)

type Location struct {
	echo       *echo.Group
	database   *sqlx.DB
	repository *database.LocationRepository
	member     *database.MemberRepository
}

func (h *Location) initialize() *Location {
	h.repository = database.CreateLocation(h.database)
	h.member = database.CreateMember(h.database)

	h.echo.POST("/locations", h.create)
	h.echo.GET("/locations", h.many)

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

func (h *Location) many(c echo.Context) error {
	org, err, status := organization(h.member, c)
	if err != nil {
		return c.JSON(status, api.Error(err.Error()))
	}

	locations, err := h.repository.ByOrganization(org)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("no locations found in the organization"))
	}

	return c.JSON(http.StatusOK, api.Success(locations, api.CreateRequest(c)))
}

func (h *Location) create(c echo.Context) error {
	location := &model.Location{}
	if err := c.Bind(location); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect location data"))
	}

	if err := c.Validate(location); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect location data", err.Error()))
	}

	claims := token.FromContext(c)
	if claims.Id == nil {
		return c.JSON(http.StatusBadRequest, api.Error("invalid authentication token"))
	}

	member, err := h.member.Find(location.OrganizationId, claims.Id)
	if err != nil || member == nil {
		return c.JSON(http.StatusForbidden, api.Error("only organization members can create locations"))
	}

	location.CreatedBy = claims.Id
	entity, err := h.repository.Create(location)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}

func CreateLocation(g *echo.Group, d *sqlx.DB) *Location {
	return (&Location{echo: g, database: d}).initialize()
}
