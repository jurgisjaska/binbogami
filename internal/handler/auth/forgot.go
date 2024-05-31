package auth

import (
	"net/http"

	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/model/auth"
	"github.com/labstack/echo/v4"
)

const passwordResetLimit int = 10

func (h *Auth) forgot(c echo.Context) error {
	request := &auth.ForgotRequest{}
	if err := c.Bind(request); err != nil {
		// request cannot be bind to model
		return c.JSON(http.StatusUnauthorized, api.Error(requestError))
	}

	if err := c.Validate(request); err != nil {
		// request does not pass validation
		return c.JSON(http.StatusUnauthorized, api.Errors(validationError, err.Error()))
	}

	// attempt to locate used by email
	user, err := h.user.FindByColumn("email", request.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.Errors("no user associated with this email", err.Error()))
	}

	resets, err := h.userPasswordReset.FindManyByUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Error("failed to retrieve reset requests"))
	}

	if len(*resets) >= passwordResetLimit {
		return c.JSON(http.StatusUnauthorized, api.Error("too many reset requests"))
	}

	request.Ip = c.RealIP()
	request.UserAgent = c.Request().UserAgent()
	request.User = user
	entity, err := h.userPasswordReset.Save(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}
