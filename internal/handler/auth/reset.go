package auth

import (
	"fmt"
	"net/http"

	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/model/auth"
	"github.com/labstack/echo/v4"
)

func (h *Auth) reset(c echo.Context) error {
	request := &auth.ResetPasswordRequest{}
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(requestError))
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, api.Errors(formError, err.Error()))
	}

	// are the passwords matching?
	if request.Password != request.RepeatedPassword {
		return c.JSON(http.StatusUnprocessableEntity, api.Errors(passwordsMatchError, fmt.Errorf("passwords does not match")))
	}

	// retrieve the password reset token
	entity, err := h.userPasswordReset.FindById(request.Token)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("password reset token not found"))
	}

	// retrieve the user that's attempting to reset password
	user, err := h.user.FindByColumn("id", entity.UserId)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.Errors(credentialError, err.Error()))
	}

	user.Password, err = h.hashPassword(request.Password, user.Salt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Errors("internal server error", err.Error()))
	}

	err = h.user.UpdatePassword(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Errors("internal server error", err.Error()))
	}

	err = h.userPasswordReset.UpdateExpireAt(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Errors("internal server error", err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(user, api.CreateRequest(c)))
}
