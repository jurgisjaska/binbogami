package public

import (
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/database/user/invitation"
	"github.com/labstack/echo/v5"
)

type Public struct {
	echo       *echo.Group
	database   *sqlx.DB
	invitation invitation.InvitationRepository
}

func (h *Public) initialize() *Public {
	h.invitation = invitation.CreateInvitation(h.database)

	h.echo.GET("/invitation/:id", h.invite)

	return h
}

func CreatePublic(e *echo.Group, d *sqlx.DB) *Public {
	return (&Public{echo: e, database: d}).initialize()
}
