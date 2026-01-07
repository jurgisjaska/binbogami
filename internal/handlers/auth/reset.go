package auth

import (
	"fmt"
	"net/http"

	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models/auth"
	"github.com/labstack/echo/v4"
)

func (h *Auth) reset(c echo.Context) error {
	request := &auth.ResetPasswordRequest{}
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(requestError))
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, api.Errors(validationError, err.Error()))
	}

	// are the passwords matching?
	if request.Password != request.RepeatedPassword {
		return c.JSON(http.StatusUnprocessableEntity, api.Errors(passwordsMatchError, fmt.Errorf("passwords does not match")))
	}

	// retrieve the password reset token
	entity, err := h.user.passwordReset.Find(request.Token)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("password reset token not found"))
	}

	// retrieve the repository that's attempting to reset password
	user, err := h.user.repository.Find(entity.UserId)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.Errors(credentialError, err.Error()))
	}

	user.Password, err = h.hashPassword(request.Password, user.Salt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Errors(internalError, err.Error()))
	}

	err = h.user.repository.UpdatePassword(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Errors(internalError, err.Error()))
	}

	err = h.user.passwordReset.UpdateExpireAt(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Errors(internalError, err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(user, api.CreateRequest(c)))
}
