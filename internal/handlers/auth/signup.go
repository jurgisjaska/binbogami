package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models/auth"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/jurgisjaska/binbogami/internal/database/user/invitation"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
)

// signup validates signup form data and creates new user
// if the invitation UUID is present and valid assigns confirmed status to the new user.
func (h *Auth) signup(c echo.Context) error {
	request := &auth.SignupRequest{}
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(requestError))
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, api.Errors(validationError, err.Error()))
	}

	if request.Password != request.RepeatedPassword {
		return c.JSON(http.StatusUnprocessableEntity, api.Errors(passwordsMatchError, fmt.Errorf("passwords does not match")))
	}

	existingUser, err := h.user.repository.FindByEmail(request.Email)
	if existingUser != nil {
		return c.JSON(http.StatusUnprocessableEntity, api.Error("email address already in use"))
	}

	role := user.RoleDefault
	var inv *invitation.Invitation
	var confirmedAt *time.Time

	if request.InvitationId != nil {
		inv, err = h.invitation.Find(*request.InvitationId)
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
		_ = h.invitation.Delete(inv)
	}

	return c.JSON(
		http.StatusOK,
		api.Success(
			auth.SignupResponse{User: u, Token: t},
			api.CreateRequest(c),
		),
	)
}
