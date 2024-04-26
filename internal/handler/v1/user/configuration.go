package user

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	um "github.com/jurgisjaska/binbogami/internal/api/model/user"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	ud "github.com/jurgisjaska/binbogami/internal/database/user"
	v1 "github.com/jurgisjaska/binbogami/internal/handler/v1"
	"github.com/labstack/echo/v4"
)

type Configuration struct {
	echo       *echo.Group
	database   *sqlx.DB
	repository *ud.ConfigurationRepository
}

func (h *Configuration) initialize() *Configuration {
	h.repository = ud.CreateConfiguration(h.database)

	h.echo.PUT("/users/configurations", h.set)

	return h
}

func (h *Configuration) set(c echo.Context) error {
	request := &um.SetConfiguration{}
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect user configuration data"))
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect user configuration data", err.Error()))
	}

	claims := token.FromContext(c)
	if claims.Id == nil {
		return c.JSON(http.StatusBadRequest, api.Error(v1.ErrorToken))
	}
	request.CreatedBy = claims.Id

	// @todo user must be a member of an organization before setting it default
	// but the organization is in the body not in the headers
	// so i need additional method to verify that

	return nil
}

func CreateConfiguration(g *echo.Group, d *sqlx.DB) *Configuration {
	return (&Configuration{echo: g, database: d}).initialize()
}
