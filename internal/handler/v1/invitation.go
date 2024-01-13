package v1

import (
	"fmt"
	"net/http"
	"net/smtp"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/model"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/database"
	"github.com/labstack/echo/v4"
)

type Invitation struct {
	echo          *echo.Group
	database      *sqlx.DB
	mail          *smtp.Client
	configuration *internal.Config
	invitation    *database.InvitationRepository
	member        *database.MemberRepository
}

func (h *Invitation) initialize() *Invitation {
	h.invitation = database.CreateInvitation(h.database)
	h.member = database.CreateMember(h.database)

	h.echo.POST("/invitations", h.create)

	return h
}

func (h *Invitation) create(c echo.Context) error {
	i := &model.Invitation{}
	if err := c.Bind(i); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("invalid invitation data"))
	}

	claims := token.FromContext(c)
	if claims.Id == nil {
		return c.JSON(http.StatusBadRequest, api.Error("invalid authentication token"))
	}

	allow := map[int]bool{
		database.MemberRoleDefault: false,
		database.MemberRoleBilling: false,
		database.MemberRoleAdmin:   true,
		database.MemberRoleOwner:   true,
	}
	cm, err := h.member.Find(i.OrganizationId, claims.Id)
	if err != nil || !allow[cm.Role] {
		return c.JSON(http.StatusForbidden, api.Error("only organization owners and admins can invite members"))
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(i); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect invitation", err.Error()))
	}

	i.CreatedBy = claims.Id
	invitations, err := h.invitation.Create(i)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	for _, invitation := range invitations {
		err = h.send(invitation)
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
		}
	}

	return c.JSON(http.StatusOK, api.Success(invitations, api.CreateRequest(c)))
}

// @todo move to separate service, the handler should not be responsible for the emails
func (h *Invitation) send(invitation *database.Invitation) error {
	if err := h.mail.Mail(h.configuration.Mail.Sender); err != nil {
		return err
	}
	if err := h.mail.Rcpt(invitation.Email); err != nil {
		return err
	}

	writer, err := h.mail.Data()
	defer func() { _ = writer.Close() }()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s:%d/join/%s", h.configuration.Hostname, h.configuration.Port, invitation.Id)
	_, err = fmt.Fprintf(writer, "Hello,\nJoin organization using the link %s", url)
	if err != nil {
		return err
	}

	return nil
}

func CreateInvitation(g *echo.Group, d *sqlx.DB, m *smtp.Client, c *internal.Config) *Invitation {
	return (&Invitation{echo: g, database: d, mail: m, configuration: c}).initialize()
}
