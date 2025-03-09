package v1

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/model"
	"github.com/jurgisjaska/binbogami/internal/database"
	"github.com/jurgisjaska/binbogami/internal/database/book"
	"github.com/jurgisjaska/binbogami/internal/database/member"
	"github.com/labstack/echo/v4"
)

type Location struct {
	echo       *echo.Group
	database   *sqlx.DB
	repository *database.LocationRepository
	member     *member.MemberRepository
	book       *book.Repository
}

func (h *Location) initialize() *Location {
	h.repository = database.CreateLocation(h.database)
	h.member = member.CreateMember(h.database)
	h.book = book.CreateBook(h.database)

	h.echo.POST("/locations", h.create)
	h.echo.GET("/locations", h.byOrganization)
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

func (h *Location) byOrganization(c echo.Context) error {
	member, err := membership(h.member, c)
	if err != nil {
		return c.JSON(http.StatusForbidden, api.Error(err.Error()))
	}

	locations, err := h.repository.ManyByOrganization(member.OrganizationId)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("no locations found in the organization"))
	}

	return c.JSON(http.StatusOK, api.Success(locations, api.CreateRequest(c)))
}

func (h *Location) create(c echo.Context) error {
	member, err := membership(h.member, c)
	if err != nil {
		return c.JSON(http.StatusForbidden, api.Error(err.Error()))
	}

	location := &model.Location{}
	if err := c.Bind(location); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect location data"))
	}

	if err := c.Validate(location); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect location data", err.Error()))
	}

	location.CreatedBy = member.UserId
	location.OrganizationId = member.OrganizationId
	entity, err := h.repository.Create(location)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}

func (h *Location) byBook(c echo.Context) error {
	_, err := membership(h.member, c)
	if err != nil {
		return c.JSON(http.StatusForbidden, api.Error(err.Error()))
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect book"))
	}

	book, err := h.book.Find(&id)
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
