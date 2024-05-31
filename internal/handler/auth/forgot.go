package auth

import (
	"net/http"

	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/model/auth"
	"github.com/labstack/echo/v4"
)

func (h *Auth) forgot(c echo.Context) error {
	request := &auth.ForgotPasswordRequest{}
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusUnauthorized, api.Error(credentialError))
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnauthorized, api.Errors(credentialError, err.Error()))
	}

	// user, err := h.user.FindByColumn("email", request.Email)
	// if err != nil {
	// 	return c.JSON(http.StatusUnauthorized, api.Errors(credentialError, err.Error()))
	// }

	// create reset token
	// send an email to user

	return nil
}
