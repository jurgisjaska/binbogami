package v1

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/database"
	"github.com/labstack/echo/v4"
)

// Organization represents an organization handler.
type Organization struct {
	echo         *echo.Group
	database     *sqlx.DB
	organization *database.OrganizationRepository
	member       *database.MemberRepository
}

func (h *Organization) initialize() *Organization {
	h.organization = database.CreateOrganization(h.database)
	h.member = database.CreateMember(h.database)

	h.echo.GET("/organizations/:id", h.one)
	// h.echo.GET("/organizations", h.many)
	h.echo.POST("/organizations", h.create)
	// h.echo.PUT("/organizations/:id", h.update)
	// h.echo.DELETE("/organizations/:id", h.delete)

	return h
}

func (h *Organization) one(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect organization"))
	}

	organization, err := h.organization.Find(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("organization not found"))
	}

	return c.JSON(http.StatusOK, api.Success(organization, api.CreateRequest(c)))
}

func (h *Organization) create(c echo.Context) error {
	// @todo this should be a api model
	organization := &database.Organization{}
	if err := c.Bind(organization); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect organization"))
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(organization); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect organization", err.Error()))
	}

	claims := token.FromContext(c)
	if claims.Id == nil {
		return c.JSON(http.StatusBadRequest, api.Error("invalid authentication token"))
	}
	organization.CreatedBy = claims.Id

	err := h.organization.Create(organization)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	// @todo create a failure recovery process
	_, _ = h.member.Create(organization.Id, claims.Id, database.MemberRoleOwner, nil)

	return c.JSON(http.StatusOK, api.Success(organization, api.CreateRequest(c)))
}

// CreateOrganization initializes and returns an instance of Organization handler.
func CreateOrganization(g *echo.Group, d *sqlx.DB) *Organization {
	return (&Organization{echo: g, database: d}).initialize()
}
