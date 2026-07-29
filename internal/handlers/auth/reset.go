package auth

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models/auth"
	"github.com/labstack/echo/v5"
)

// open sets the opened_at field to the current time.
func (h *Auth) open(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.auditlog.Warn("password reset open error: incorrect password reset token", "error", err.Error())
		return c.JSON(http.StatusBadRequest, api.Error("incorrect password reset token"))
	}

	// retrieve the password reset token
	entity, err := h.user.passwordReset.Find(id)
	if err != nil {
		h.auditlog.Warn("password reset open error: password reset token not found", "error", err.Error())
		return c.JSON(http.StatusNotFound, api.Error("password reset token not found"))
	}

	n := time.Now()
	entity.OpenedAt = &n
	err = h.user.passwordReset.Update(entity)
	if err != nil {
		h.auditlog.Warn("password reset open error: failed to update token", "error", err.Error())
		return c.JSON(http.StatusInternalServerError, api.Error("failed to update password reset token"))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}

func (h *Auth) reset(c *echo.Context) error {
	request := &auth.ResetPasswordRequest{}
	if err := c.Bind(request); err != nil {
		h.auditlog.Warn("password reset error: bad request", "error", err.Error())
		return c.JSON(http.StatusBadRequest, api.Error(requestError))
	}

	if err := c.Validate(request); err != nil {
		h.auditlog.Warn("password reset error: validation error", "error", err.Error())
		return c.JSON(http.StatusUnprocessableEntity, api.Errors(validationError, err.Error()))
	}

	// retrieve the password reset token
	entity, err := h.user.passwordReset.Find(request.Token)
	if err != nil {
		h.auditlog.Warn("password reset error: password reset token not found", "error", err.Error(), "token", request.Token)
		return c.JSON(http.StatusNotFound, api.Error("password reset token not found"))
	}

	// check if the token is opened
	// if token is not opened, password reset is not allowed
	if entity.OpenedAt == nil {
		h.auditlog.Warn("password reset error: password reset token not opened", "token", request.Token)
		return c.JSON(http.StatusUnauthorized, api.Error("password reset token not opened"))
	}

	// retrieve the user attempting to reset password
	user, err := h.user.repository.Find(entity.UserId)
	if err != nil {
		h.auditlog.Warn("password reset error: user not found", "error", err.Error(), "token", request.Token)
		return c.JSON(http.StatusUnauthorized, api.Errors(credentialError, err.Error()))
	}

	user.Password, err = h.hashPassword(request.Password, user.Salt)
	if err != nil {
		h.auditlog.Warn("password reset error: password hashing failure", "error", err.Error(), "token", request.Token)
		return c.JSON(http.StatusInternalServerError, api.Errors(internalError, err.Error()))
	}

	err = h.user.repository.UpdatePassword(user)
	if err != nil {
		h.auditlog.Warn("password reset error: failed to update password", "error", err.Error(), "token", request.Token)
		return c.JSON(http.StatusInternalServerError, api.Errors(internalError, err.Error()))
	}

	n := time.Now()
	entity.ExpireAt = n
	_ = h.user.passwordReset.Update(entity)

	return c.JSON(http.StatusOK, api.Success(user, api.CreateRequest(c)))
}
