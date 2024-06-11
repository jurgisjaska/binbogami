package auth

import (
	"fmt"
	"net/http"

	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/model/auth"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// signin in creates new JWT token for the user if credentials are correct
func (h *Auth) signin(c echo.Context) error {
	request := &auth.SigninRequest{}
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusUnauthorized, api.Error(credentialError))
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnauthorized, api.Errors(credentialError, err.Error()))
	}

	u, err := h.userRepository.user.FindByColumn("email", request.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.Errors(credentialError, err.Error()))
	}

	password := fmt.Sprintf("%s%s%s", request.Password, u.Salt, h.configuration.Secret)
	if len(password) > 71 {
		password = password[:71]
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return c.JSON(http.StatusUnauthorized, api.Error(err.Error()))
	}

	t, err := token.CreateToken(u, h.configuration.Secret)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.Error(err.Error()))
	}

	m, o := h.membership(u)
	response := auth.SigninResponse{Token: t, User: u, Member: m, Organization: o}

	return c.JSON(http.StatusOK, api.Success(response, api.CreateRequest(c)))
}
