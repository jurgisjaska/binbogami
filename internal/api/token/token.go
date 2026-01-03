package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// expire is a constant representing the number of hours for token expiration.
const expire = 72

// Claims represents the claims contained in a JWT token.
type Claims struct {
	Id    *uuid.UUID `json:"id"`
	Email *string    `json:"email"`
	Name  string     `json:"name"`
	jwt.RegisteredClaims
}

// CreateToken creates a JWT token string for a given user.
func CreateToken(u *user.User, key string) (string, error) {
	expire := jwt.NewNumericDate(time.Now().Add(time.Hour * expire))
	claim := &Claims{
		Id:    &u.Id,
		Email: &u.Email,
		Name:  fmt.Sprintf("%s %s", u.Name, u.Surname),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expire,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(key))
}

// CreateJWTConfig creates a JWT token configuration for Echo framework.
// The configuration includes a function to create new claims and the signing key.
func CreateJWTConfig(key string) echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(Claims)
		},
		SigningKey: []byte(key),
	}
}

func FromContext(c echo.Context) *Claims {
	token := c.Get("user").(*jwt.Token)
	return token.Claims.(*Claims)
}
