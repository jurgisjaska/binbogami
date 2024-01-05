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
	echo       *echo.Group
	database   *sqlx.DB
	repository *database.OrganizationRepository
}

func (h *Organization) initialize() *Organization {
	h.repository = database.CreateOrganization(h.database)
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

	organization, err := h.repository.Find(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("organization not found"))
	}

	return c.JSON(http.StatusOK, api.Success(organization, api.CreateRequest(c)))
}

func (h *Organization) create(c echo.Context) error {
	organization := &database.Organization{}
	if err := c.Bind(organization); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect organization"))
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(organization); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect organization", err.Error()))
	}

	// @todo single user should not be able to create multiple organization with the same name
	// @todo single user should not be able to own multiple organizations with the same name

	claims := token.FromContext(c)
	if claims.Id == nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect credentials"))
	}
	organization.CreatedBy = claims.Id
	organization.OwnedBy = claims.Id

	err := h.repository.Create(organization)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(organization, api.CreateRequest(c)))
}

// CreateOrganization initializes and returns an instance of Organization handler.
func CreateOrganization(g *echo.Group, d *sqlx.DB) *Organization {
	return (&Organization{echo: g, database: d}).initialize()
}
