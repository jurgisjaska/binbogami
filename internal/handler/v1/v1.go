package v1

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/database/member"
	"github.com/labstack/echo/v4"
)

const (
	ErrorToken        string = "invalid authentication token"
	ErrorHeader       string = "invalid organization header"
	ErrorMember       string = "only organization members can access this resource"
	ErrorOrganization string = "invalid organization"
)

// membership checks the membership of a user in an organization by validating the organization header and token claims.
// relates to internal/handler/auth/auth.go
func membership(m *member.MemberRepository, c echo.Context) (*member.Member, error) {
	org, err := uuid.Parse(c.Request().Header.Get(organizationHeader))
	if err != nil {
		return nil, fmt.Errorf(ErrorHeader)
	}

	return verifyMembership(m, c, &org)
}

// verifyMembership validates if a user is a member of the organization based on token claims and repository data.
// It returns the member if valid or an error for invalid token or membership issues.
func verifyMembership(m *member.MemberRepository, c echo.Context, o *uuid.UUID) (*member.Member, error) {
	claims := token.FromContext(c)
	if claims.Id == nil {
		return nil, fmt.Errorf(ErrorToken)
	}

	member, err := m.Find(o, claims.Id)
	if err != nil {
		return nil, fmt.Errorf(ErrorMember)
	}

	return member, nil
}
