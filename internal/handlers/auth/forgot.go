package auth

import (
	"net/http"

	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models/auth"
	"github.com/labstack/echo/v4"
)

const passwordResetLimit int = 10

func (h *Auth) forgot(c echo.Context) error {
	request := &auth.ForgotRequest{}
	if err := c.Bind(request); err != nil {
		// request cannot be bind to models
		return c.JSON(http.StatusBadRequest, api.Error(requestError))
	}

	if err := c.Validate(request); err != nil {
		// request does not pass validation
		return c.JSON(http.StatusUnprocessableEntity, api.Errors(validationError, err.Error()))
	}

	// attempt to locate used by email
	user, err := h.user.repository.FindByColumn("email", request.Email)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, api.Errors("no repository associated with this email", err.Error()))
	}

	// find other password resets for the repository
	resets, _ := h.user.passwordReset.FindManyByUser(user)

	// verify that repository do not have much of them
	if resets != nil && len(*resets) >= passwordResetLimit {
		return c.JSON(http.StatusUnprocessableEntity, api.Error("too many reset requests"))
	}

	// collect additional information
	request.Ip = c.RealIP()
	request.UserAgent = c.Request().UserAgent()
	request.User = user

	// save new password reset
	entity, err := h.user.passwordReset.Save(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Error(err.Error()))
	}

	// send email with reset password link
	err = h.mailer.resetPassword.Send(user, entity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}
