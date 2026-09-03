package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/labstack/echo/v5"
)

// UserActive returns an Echo middleware that checks if the user from the JWT token is active
// (i.e., not deleted and confirmed). If the user is deleted or unconfirmed, the request is denied with 401 Unauthorized.
func UserActive(repo user.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			v := c.Get("user")
			if v == nil {
				return c.JSON(http.StatusUnauthorized, api.Error("unauthorized"))
			}

			t, ok := v.(*jwt.Token)
			if !ok || t == nil {
				return c.JSON(http.StatusUnauthorized, api.Error("unauthorized"))
			}

			claims, ok := t.Claims.(*token.Claims)
			if !ok || claims == nil || claims.Id == nil {
				return c.JSON(http.StatusUnauthorized, api.Error("unauthorized"))
			}

			u, err := repo.Find(*claims.Id)
			if err != nil || u == nil {
				return c.JSON(http.StatusUnauthorized, api.Error("user not found"))
			}

			if u.DeletedAt != nil {
				return c.JSON(http.StatusUnauthorized, api.Error("user deleted"))
			}

			if u.ConfirmedAt == nil {
				return c.JSON(http.StatusUnauthorized, api.Error("user not confirmed"))
			}

			return next(c)
		}
	}
}

// CreateUserActive creates a UserActive middleware for Echo.
func CreateUserActive(repo user.UserRepository) echo.MiddlewareFunc {
	return UserActive(repo)
}
