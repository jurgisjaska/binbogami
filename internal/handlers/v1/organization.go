package v1

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/database/member"
	"github.com/jurgisjaska/binbogami/internal/database/organization"
	"github.com/labstack/echo/v4"
)

const (
	organizationHeader = "organization"
)

// Organization represents an organization handlers.
type Organization struct {
	echo       *echo.Group
	database   *sqlx.DB
	repository *organization.Repository
	member     *member.MemberRepository
}

func (h *Organization) initialize() *Organization {
	h.repository = organization.CreateOrganization(h.database)
	h.member = member.CreateMember(h.database)

	h.echo.GET("/organizations/:id", h.show)
	h.echo.GET("/organizations", h.byMember)
	h.echo.POST("/organizations", h.create)

	return h
}

func (h *Organization) show(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(ErrorOrganization))
	}

	claims := token.FromContext(c)
	if claims.Id == nil {
		return c.JSON(http.StatusBadRequest, api.Error(ErrorToken))
	}

	organization, err := h.repository.Find(&id, claims.Id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("no organizations found"))
	}

	return c.JSON(http.StatusOK, api.Success(organization, api.CreateRequest(c)))
}

func (h *Organization) create(c echo.Context) error {
	organization := &models.Organization{}
	if err := c.Bind(organization); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect organization"))
	}

	if err := c.Validate(organization); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect organization", err.Error()))
	}

	claims := token.FromContext(c)
	if claims.Id == nil {
		return c.JSON(http.StatusBadRequest, api.Error(ErrorToken))
	}
	organization.CreatedBy = claims.Id

	_, err := h.repository.ByMemberAndName(claims.Id, organization.Name)
	if err == nil {
		return c.JSON(http.StatusBadRequest, api.Error("organization with the same name already exists"))
	}

	entity, err := h.repository.Create(organization)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	// @todo create a failure recovery process
	_, _ = h.member.Create(entity.Id, claims.Id, member.MemberRoleOwner, claims.Id)

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}

// byMember handles the GET request to retrieve organizations where the user is a member.
func (h *Organization) byMember(c echo.Context) error {
	claims := token.FromContext(c)
	if claims.Id == nil {
		return c.JSON(http.StatusBadRequest, api.Error(ErrorToken))
	}

	organizations, err := h.repository.ByMember(claims.Id)

	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("no organizations found"))
	}

	return c.JSON(http.StatusOK, api.Success(organizations, api.CreateRequest(c)))
}

// CreateOrganization initializes and returns an instance of Organization handlers.
func CreateOrganization(g *echo.Group, d *sqlx.DB) *Organization {
	return (&Organization{echo: g, database: d}).initialize()
}
