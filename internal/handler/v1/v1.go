package v1

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/database"
	"github.com/labstack/echo/v4"
)

const (
	errorToken  string = "invalid authentication token"
	errorHeader string = "invalid organization header"
	errorMember string = "only organization members can access this resource"
)

// organization retrieves the organization ID from the request header
// and performs authorization checks based on the given member repository.
// @todo update all endpoints related with the organization tu use the same header
func organization(m *database.MemberRepository, c echo.Context) (*uuid.UUID, error) {
	org, err := uuid.Parse(c.Request().Header.Get(organizationHeader))
	if err != nil {
		return nil, fmt.Errorf(errorHeader)
	}

	claims := token.FromContext(c)
	if claims.Id == nil {
		return nil, fmt.Errorf(errorToken)
	}

	member, err := m.Find(&org, claims.Id)
	if err != nil || member == nil {
		return nil, fmt.Errorf(errorMember)
	}

	return &org, nil
}
