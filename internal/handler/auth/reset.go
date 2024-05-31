package auth

import (
	"net/http"

	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/model/auth"
	"github.com/labstack/echo/v4"
)

func (h *Auth) reset(c echo.Context) error {
	request := &auth.ResetPasswordRequest{}
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusUnauthorized, api.Error(credentialError))
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnauthorized, api.Errors(credentialError, err.Error()))
	}

	// unique reset token is needed to reset password
	// it can be used once
	// expires after 24 hours

	return nil
}
