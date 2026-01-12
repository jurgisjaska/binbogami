package public

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/labstack/echo/v4"
)

func (h *Public) reset(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect password reset token"))
	}

	entity, err := h.passwordReset.Find(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("password reset token not found"))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}
