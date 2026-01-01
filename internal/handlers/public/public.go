package public

import (
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/database/organization"
	"github.com/jurgisjaska/binbogami/internal/database/user/invitation"
	"github.com/jurgisjaska/binbogami/internal/database/user/password"
	"github.com/labstack/echo/v4"
)

type Public struct {
	echo          *echo.Group
	database      *sqlx.DB
	invitation    *invitation.InvitationRepository
	passwordReset *password.PasswordResetRepository
}

func (h *Public) initialize() *Public {
	h.invitation = invitation.CreateInvitation(h.database)
	h.passwordReset = password.CreatePasswordReset(h.database)

	h.echo.GET("/invitation/:id", h.invite)
	h.echo.GET("/reset-password/:id", h.reset)

	return h
}

func CreatePublic(e *echo.Group, d *sqlx.DB) *Public {
	return (&Public{echo: e, database: d}).initialize()
}
