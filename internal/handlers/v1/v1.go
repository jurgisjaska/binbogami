package v1

import (
	"fmt"

	"github.com/jurgisjaska/binbogami/internal/api/token"
	u "github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/labstack/echo/v4"
)

const (
	ErrorToken string = "invalid authentication token"
)

func currentUser(r *u.Repository, c echo.Context) (*u.User, error) {
	claims := token.FromContext(c)
	if claims.Id == nil {
		return nil, fmt.Errorf(ErrorToken)
	}

	return r.FindActive(*claims.Id)
}
