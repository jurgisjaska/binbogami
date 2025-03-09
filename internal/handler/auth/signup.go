package auth

import (
	"fmt"
	"net/http"

	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/model/auth"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/database/member"
	"github.com/jurgisjaska/binbogami/internal/database/organization"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/labstack/gommon/random"
)

// signup validates signup form data and creates new user
// if the invitation UUID is present adds the new user to the organization
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

	existingUser, err := h.userRepository.user.FindByColumn("email", *request.Email)
	if existingUser != nil {
		return c.JSON(http.StatusUnprocessableEntity, api.Error("email address already in use"))
	}

	u := &user.User{
		Email:   request.Email,
		Name:    request.Name,
		Surname: request.Surname,
		Salt:    random.String(16),
	}

	u.Password, err = h.hashPassword(request.Password, u.Salt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Errors(internalError, err.Error()))
	}

	err = h.userRepository.user.Create(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Errors(internalError, err.Error()))
	}

	t, err := token.CreateToken(u, h.configuration.Secret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Error(err.Error()))
	}

	// @todo should invitation and user email match when using invitation link?

	member := &member.Member{}
	organization := &organization.Organization{}
	log.Infof("%+v", request)
	if request.InvitationId != nil {
		invitation, err := h.invitation.FindById(request.InvitationId)
		log.Infof("%+v", invitation)
		log.Error(err)
		if err == nil {
			member, err = h.member.Create(invitation.OrganizationId, u.Id, member.MemberRoleDefault, invitation.CreatedBy)
			if err == nil {
				_ = h.invitation.Delete(invitation)
			}

			// error does not matter in this case
			// organization either is there or no
			// SQL not found can be ignored
			organization, _ = h.organization.FindById(invitation.OrganizationId)
		}
	}

	// membership status
	m := false
	if member.Id != 0 {
		m = true
	}

	// reset organization to nil
	// to keep consistency between sign in and sign up methods
	if organization.Id == nil {
		organization = nil
	}

	return c.JSON(
		http.StatusOK,
		api.Success(
			auth.SignupResponse{User: u, Token: t, Member: m, Organization: organization},
			api.CreateRequest(c),
		),
	)
}
