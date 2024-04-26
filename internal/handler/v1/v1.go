package v1

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/database"
	"github.com/labstack/echo/v4"
)

const (
	ErrorToken        string = "invalid authentication token"
	ErrorHeader       string = "invalid organization header"
	ErrorMember       string = "only organization members can access this resource"
	ErrorOrganization string = "invalid organization"
)

// membership checks the membership of a user in an organization by validating the organization header and token claims.
func membership(m *database.MemberRepository, c echo.Context) (*database.Member, error) {
	org, err := uuid.Parse(c.Request().Header.Get(organizationHeader))
	if err != nil {
		return nil, fmt.Errorf(ErrorHeader)
	}

	claims := token.FromContext(c)
	if claims.Id == nil {
		return nil, fmt.Errorf(ErrorToken)
	}

	member, err := m.Find(&org, claims.Id)
	if err != nil {
		return nil, fmt.Errorf(ErrorMember)
	}

	return member, nil
}
