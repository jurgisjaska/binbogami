package auth

import (
	"net/http"

	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models/auth"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/bcrypt"
)

// signin in creates a new JWT token for the user if credentials are correct
func (h *Auth) signin(c *echo.Context) error {
	request := &auth.SigninRequest{}
	if err := c.Bind(request); err != nil {
		h.auditlog.Warn("signin error", "error", err.Error())
		return c.JSON(http.StatusUnauthorized, api.Error(credentialError))
	}

	if err := c.Validate(request); err != nil {
		h.auditlog.Warn("signin error", "error", err.Error())
		return c.JSON(http.StatusUnauthorized, api.Errors(credentialError, err.Error()))
	}

	u, err := h.user.repository.FindActiveByEmail(request.Email)
	if err != nil {
		h.auditlog.Warn("signin error", "email", request.Email, "error", err.Error())
		return c.JSON(http.StatusUnauthorized, api.Errors(credentialError, err.Error()))
	}

	password := h.buildPassword(request.Password, u.Salt)
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		h.auditlog.Warn("signin error", "user_id", u.Id, "error", err.Error())
		return c.JSON(http.StatusUnauthorized, api.Error(err.Error()))
	}

	t, err := token.CreateToken(u, h.configuration.Secret)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.Error(err.Error()))
	}

	response := auth.SigninResponse{Token: t, User: u}
	h.auditlog.Info("signed in", "user_id", u.Id)

	return c.JSON(http.StatusOK, api.Success(response, api.CreateRequest(c)))
}
