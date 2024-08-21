package v1

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/model"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/database"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/jurgisjaska/binbogami/internal/service/mail"
	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
)

type Invitation struct {
	echo          *echo.Group
	database      *sqlx.DB
	mailer        *mail.Invitation
	configuration *internal.Config
	invitation    *database.InvitationRepository
	member        *database.MemberRepository
	organization  *database.OrganizationRepository
	user          *user.Repository
}

func (h *Invitation) initialize() *Invitation {
	h.invitation = database.CreateInvitation(h.database)
	h.member = database.CreateMember(h.database)
	h.organization = database.CreateOrganization(h.database)
	h.user = user.CreateUser(h.database)

	h.echo.POST("/invitations", h.create)
	h.echo.GET("/invitations", h.byOrganizationMember)

	return h
}

func (h *Invitation) byOrganizationMember(c echo.Context) error {
	member, err := membership(h.member, c)
	if err != nil {
		return c.JSON(http.StatusForbidden, api.Error(err.Error()))
	}

	invitations, err := h.invitation.FindByMember(member)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("no invitations found for current organization"))
	}

	return c.JSON(http.StatusOK, api.Success(invitations, api.CreateRequest(c)))
}

// if request data are strings
// failure happens only with mailer
// this results in bunch of incorrect data in database
func (h *Invitation) create(c echo.Context) error {
	invitation := &model.InvitationRequest{}
	if err := c.Bind(invitation); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("invalid invitation"))
	}

	claims := token.FromContext(c)
	if claims.Id == nil {
		return c.JSON(http.StatusUnauthorized, api.Error(ErrorToken))
	}

	allow := map[int]bool{
		database.MemberRoleDefault: false,
		database.MemberRoleBilling: false,
		database.MemberRoleAdmin:   true,
		database.MemberRoleOwner:   true,
	}

	member, err := membership(h.member, c)
	if err != nil || !allow[member.Role] {
		return c.JSON(http.StatusForbidden, api.Error("only organization owners and admins can invite members"))
	}

	if err = c.Validate(invitation); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, api.Errors("incorrect invitation", err.Error()))
	}

	invitation.CreatedBy = claims.Id
	invitation.OrganizationId = member.OrganizationId
	invitations, err := h.invitation.Create(invitation)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	organization, err := h.organization.FindById(invitation.OrganizationId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Error(err.Error()))
	}

	sender, err := h.user.FindByColumn("id", invitation.CreatedBy)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Error(err.Error()))
	}

	for _, i := range invitations {
		err = h.mailer.Send(sender, organization, i)
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
		}
	}

	return c.JSON(http.StatusOK, api.Success(invitations, api.CreateRequest(c)))
}

// CreateInvitation creates a new Invitation handler and initializes it.
func CreateInvitation(g *echo.Group, d *sqlx.DB, c *internal.Config, md *gomail.Dialer) *Invitation {
	return (&Invitation{
		echo:          g,
		database:      d,
		configuration: c,
		mailer:        mail.CreateInvitation(md, c),
	}).initialize()
}
