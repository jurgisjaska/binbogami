package user

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	um "github.com/jurgisjaska/binbogami/internal/api/model/user"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/database"
	ud "github.com/jurgisjaska/binbogami/internal/database/user"
	v1 "github.com/jurgisjaska/binbogami/internal/handler/v1"
	"github.com/labstack/echo/v4"
)

type Configuration struct {
	echo         *echo.Group
	database     *sqlx.DB
	member       *database.MemberRepository
	organization *database.OrganizationRepository
	repository   *ud.ConfigurationRepository
}

func (h *Configuration) initialize() *Configuration {
	h.repository = ud.CreateConfiguration(h.database)
	h.member = database.CreateMember(h.database)
	h.organization = database.CreateOrganization(h.database)

	h.echo.PUT("/users/configurations", h.set)

	return h
}

func (h *Configuration) set(c echo.Context) error {
	request := &um.SetConfigurationRequest{}
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect user configuration data", err.Error()))
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect user configuration data", err.Error()))
	}

	claims := token.FromContext(c)
	if claims.Id == nil {
		return c.JSON(http.StatusBadRequest, api.Error(v1.ErrorToken))
	}
	request.UserId = claims.Id

	organization := uuid.MustParse(request.Value)
	_, err := h.member.Find(&organization, claims.Id)
	if err != nil {
		return c.JSON(http.StatusForbidden, api.Error(err.Error()))
	}

	entity, err := h.repository.Upsert(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	o, err := h.organization.FindById(&organization)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("organization not found"))
	}

	response := &um.ConfigurationResponse{entity, o}
	return c.JSON(http.StatusOK, api.Success(response, api.CreateRequest(c)))
}

func CreateConfiguration(g *echo.Group, d *sqlx.DB) *Configuration {
	return (&Configuration{echo: g, database: d}).initialize()
}
