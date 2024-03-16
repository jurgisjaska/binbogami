package p

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/model"
	"github.com/jurgisjaska/binbogami/internal/database"
	"github.com/labstack/echo/v4"
)

// Invitation represents an invitation handler.
type Invitation struct {
	echo         *echo.Group
	database     *sqlx.DB
	repository   *database.InvitationRepository
	organization *database.OrganizationRepository
}

func (h *Invitation) initialize() *Invitation {
	h.repository = database.CreateInvitation(h.database)
	h.organization = database.CreateOrganization(h.database)

	h.echo.GET("/invitation/:id", h.invitation)

	return h
}

func (h *Invitation) invitation(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect invitation"))
	}

	invitation, err := h.repository.Open(&id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("invitation not found"))
	}

	organization, err := h.organization.ById(invitation.OrganizationId)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("organization not found"))
	}

	response := &model.InvitationResponse{invitation, organization}
	return c.JSON(http.StatusOK, api.Success(response, api.CreateRequest(c)))
}

// CreateInvitation initializes a new Invitation object.
func CreateInvitation(e *echo.Group, d *sqlx.DB) *Invitation {
	return (&Invitation{echo: e, database: d}).initialize()
}
