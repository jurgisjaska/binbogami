package p

import (
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/database"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/labstack/echo/v4"
)

type Public struct {
	echo          *echo.Group
	database      *sqlx.DB
	invitation    *database.InvitationRepository
	organization  *database.OrganizationRepository
	passwordReset *user.PasswordResetRepository
}

func (h *Public) initialize() *Public {
	h.invitation = database.CreateInvitation(h.database)
	h.organization = database.CreateOrganization(h.database)
	h.passwordReset = user.CreatePasswordReset(h.database)

	h.echo.GET("/invitation/:id", h.invite)
	h.echo.GET("/reset-password/:id", h.reset)

	return h
}

func CreatePublic(e *echo.Group, d *sqlx.DB) *Public {
	return (&Public{echo: e, database: d}).initialize()
}
