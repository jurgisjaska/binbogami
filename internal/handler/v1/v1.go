package v1

import (
	"fmt"
	"net/http"

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
func organization(m *database.MemberRepository, c echo.Context) (*uuid.UUID, error, int) {
	org, err := uuid.Parse(c.Request().Header.Get(organizationHeader))
	if err != nil {
		return nil, fmt.Errorf(errorHeader), http.StatusBadRequest
	}

	claims := token.FromContext(c)
	if claims.Id == nil {
		return nil, fmt.Errorf(errorToken), http.StatusBadRequest
	}

	member, err := m.Find(&org, claims.Id)
	if err != nil || member == nil {
		return nil, fmt.Errorf(errorMember), http.StatusBadRequest
	}

	return &org, nil, http.StatusOK
}
