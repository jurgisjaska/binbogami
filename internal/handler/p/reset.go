package p

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/labstack/echo/v4"
)

// Reset represents a type for getting public information about password reset.
type Reset struct {
	echo       *echo.Group
	database   *sqlx.DB
	repository *user.PasswordResetRepository
}

func (h *Reset) initialize() *Reset {
	h.repository = user.CreatePasswordReset(h.database)

	h.echo.GET("/password-reset/:id", h.passwordReset)

	return h
}

func (h *Reset) passwordReset(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect invitation"))
	}

	entity, err := h.repository.FindById(&id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("password reset token not found"))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}

func CreateReset(e *echo.Group, d *sqlx.DB) *Reset {
	return (&Reset{echo: e, database: d}).initialize()
}
