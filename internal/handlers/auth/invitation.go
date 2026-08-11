package auth

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/labstack/echo/v5"
)

func (h *Auth) openInvitation(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.auditlog.Warn("invitation open error: incorrect invitation", "error", err.Error())
		return c.JSON(http.StatusBadRequest, api.Error("incorrect invitation"))
	}

	invitation, err := h.invitation.Find(id)
	if err != nil {
		h.auditlog.Warn("invitation open error: invitation not found", "error", err.Error())
		return c.JSON(http.StatusNotFound, api.Error("invitation not found"))
	}

	if invitation.OpenedAt == nil {
		n := time.Now()
		invitation.OpenedAt = &n
		err = h.invitation.Update(invitation)
		if err != nil {
			h.auditlog.Warn("invitation open error: failed to update invitation", "error", err.Error())
			return c.JSON(http.StatusInternalServerError, api.Error("failed to update invitation"))
		}
	}

	h.auditlog.Info("invitation opened", "invitation_id", id)
	return c.JSON(http.StatusOK, api.Success(invitation, api.CreateRequest(c)))
}
