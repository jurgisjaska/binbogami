package auth

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models/auth"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/jurgisjaska/binbogami/internal/database/user/invitation"
	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/random"
)

// signup validates signup form data and creates a new user
// if the invitation UUID is present and valid assigns confirmed status to the new user.
func (h *Auth) signup(c *echo.Context) error {
	request := &auth.SignupRequest{}
	if err := c.Bind(request); err != nil {
		h.auditlog.Warn("signup error: bad request", "error", err.Error())
		return c.JSON(http.StatusBadRequest, api.Error(requestError))
	}

	if err := c.Validate(request); err != nil {
		h.auditlog.Warn("signup error: validation error", "error", err.Error())
		return c.JSON(http.StatusUnprocessableEntity, api.Errors(validationError, err.Error()))
	}

	existingUser, err := h.user.repository.FindByEmail(request.Email)
	if existingUser != nil {
		h.auditlog.Warn("signup error: bad request", "email", request.Email)
		return c.JSON(http.StatusUnprocessableEntity, api.Error("email address already in use"))
	}

	role := user.RoleDefault
	var inv *invitation.Invitation
	var confirmedAt *time.Time

	if request.Invitation != nil {
		inv, err = h.invitation.Find(*request.Invitation)
		if err == nil {
			n := time.Now()
			confirmedAt = &n

			if inv.Role != nil {
				role = *inv.Role
			}
		}
	}

	u := &user.User{
		Id:          uuid.New(),
		Email:       request.Email,
		Name:        request.Name,
		Surname:     request.Surname,
		Position:    request.Position,
		Salt:        random.String(16),
		Role:        role,
		CreatedAt:   time.Now(),
		ConfirmedAt: confirmedAt,
	}

	u.Password, err = h.hashPassword(request.Password, u.Salt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Errors(internalError, err.Error()))
	}

	err = h.user.repository.Create(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Errors(internalError, err.Error()))
	}

	t, err := token.CreateToken(u, h.configuration.Secret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Error(err.Error()))
	}

	if inv != nil {
		inv.UserId = &u.Id
		_ = h.invitation.Delete(inv)
	}

	// @todo: send welcome email
	// @todo: send an email confirm request

	return c.JSON(
		http.StatusOK,
		api.Success(
			auth.SignupResponse{User: u, Token: t},
			api.CreateRequest(c),
		),
	)
}
