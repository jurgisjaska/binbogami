package p

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models"
	"github.com/labstack/echo/v4"
)

func (h *Public) invite(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect invitation"))
	}

	invitation, err := h.invitation.Open(&id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("invitation not found"))
	}

	organization, err := h.organization.FindById(invitation.OrganizationId)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("organization not found"))
	}

	response := &models.InvitationResponse{Invitation: invitation, Organization: organization}
	return c.JSON(http.StatusOK, api.Success(response, api.CreateRequest(c)))
}
