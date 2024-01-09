package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/database"
	"github.com/labstack/echo/v4"
)

type (
	Public struct {
		echo                 *echo.Echo
		database             *sqlx.DB
		invitationRepository *database.InvitationRepository
	}
)

func (h *Public) initialize() *Public {
	h.invitationRepository = database.CreateInvitation(h.database)

	h.echo.GET("/join/:id", h.join)

	return h
}

func (h *Public) join(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect invitation"))
	}

	invitation, err := h.invitationRepository.Open(&id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("invitation not found"))
	}

	return c.JSON(http.StatusOK, api.Success(invitation, api.CreateRequest(c)))
}

func CreatePublic(e *echo.Echo, d *sqlx.DB) *Public {
	return (&Public{echo: e, database: d}).initialize()
}
