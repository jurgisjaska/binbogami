package auth

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models/auth"
	"github.com/jurgisjaska/binbogami/internal/database/user/password"
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

	// attempt to locate user (active) by email
	user, err := h.user.repository.FindActiveByEmail(request.Email)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, api.Errors("no repository associated with this email", err.Error()))
	}

	// find other password resets for the repository
	resets, err := h.user.passwordReset.FindManyByUser(user, passwordResetLimit)

	// verify that repository does not have much of them
	if resets != nil && len(*resets) >= passwordResetLimit {
		return c.JSON(http.StatusUnprocessableEntity, api.Error("too many reset requests"))
	}

	reset := &password.Reset{
		Id:        uuid.New(),
		UserId:    user.Id,
		Ip:        c.RealIP(),
		UserAgent: c.Request().UserAgent(),
		CreatedAt: time.Now(),
		ExpireAt:  time.Now().Add(time.Hour * password.DefaultPasswordResetDuration),
	}

	// save new password reset
	err = h.user.passwordReset.Create(reset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Error(err.Error()))
	}

	// send email with reset password link
	err = h.mailer.resetPassword.Send(user, reset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(reset, api.CreateRequest(c)))
}
