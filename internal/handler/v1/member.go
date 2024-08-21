package v1

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/model"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/database"
	"github.com/labstack/echo/v4"
)

type Member struct {
	echo       *echo.Group
	database   *sqlx.DB
	repository *database.MemberRepository
}

func (h *Member) initialize() *Member {
	h.repository = database.CreateMember(h.database)

	h.echo.GET("/members", h.one)
	h.echo.POST("/members", h.create)

	return h
}

func (h *Member) one(c echo.Context) error {
	member, err := membership(h.repository, c)
	if err != nil {
		return c.JSON(http.StatusForbidden, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(member, api.CreateRequest(c)))
}

// create adds new member to the organization
func (h *Member) create(c echo.Context) error {
	m := &model.Member{}
	if err := c.Bind(m); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("invalid member data"))
	}

	claims := token.FromContext(c)
	if claims.Id == nil {
		return c.JSON(http.StatusBadRequest, api.Error(ErrorToken))
	}

	allow := map[int]bool{
		database.MemberRoleDefault: false,
		database.MemberRoleBilling: false,
		database.MemberRoleAdmin:   true,
		database.MemberRoleOwner:   true,
	}
	cm, err := h.repository.Find(m.OrganizationId, claims.Id)
	if err != nil || !allow[cm.Role] {
		return c.JSON(http.StatusForbidden, api.Error("only organization owners and admins can add members"))
	}

	// the organization can have only single member who is an owner
	// for that reason it should be impossible to created member with role owner
	if m.Role == database.MemberRoleOwner {
		return c.JSON(http.StatusConflict, api.Error("organization already have an owner"))
	}

	member, err := h.repository.Create(m.OrganizationId, m.UserId, m.Role, claims.Id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(member, api.CreateRequest(c)))
}

func CreateMember(g *echo.Group, d *sqlx.DB) *Member {
	return (&Member{echo: g, database: d}).initialize()
}
