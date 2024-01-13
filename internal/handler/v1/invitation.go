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
	repository    *database.InvitationRepository
}

func (h *Invitation) initialize() *Invitation {
	h.repository = database.CreateInvitation(h.database)
	// h.echo.GET("/invitations/:id", h.one)
	// h.echo.GET("/invitations", h.many)
	h.echo.POST("/invitations", h.create)
	// h.echo.PUT("/invitations/:id", h.update)
	// h.echo.DELETE("/invitations/:id", h.delete)

	return h
}

func (h *Invitation) create(c echo.Context) error {
	invitation := &model.Invitation{}
	if err := c.Bind(invitation); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect invitation request"))
	}

	// @todo find a better way to this instead of repeating
	claims := token.FromContext(c)
	if claims.Id == nil {
		return c.JSON(http.StatusBadRequest, api.Error("invalid authentication token"))
	}

	// @todo validation should be done during .Bind (https://github.com/labstack/echo/issues/438)
	// @todo organization validation: exists and author; right now there will be database error but that is not enough
	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(invitation); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect invitation", err.Error()))
	}

	invitation.Author = claims.Id
	invitations, err := h.repository.Create(invitation)
	if err != nil {
		// @todo need to cleanup and make proper error handling as right now it's a mess
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
