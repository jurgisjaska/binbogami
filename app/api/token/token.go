package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/app/database"
)

const (
	expire = 72
)

type (
	Claims struct {
		Id    *uuid.UUID `json:"id"`
		Email *string    `json:"email"`
		Name  *string    `json:"name"`
		jwt.RegisteredClaims
	}
)

// CreateToken creates new JTW token instance from user model
func CreateToken(u *database.User, key string) (string, error) {
	expire := jwt.NewNumericDate(time.Now().Add(time.Hour * expire))
	claim := &Claims{
		Id:    u.Id,
		Email: u.Email,
		Name:  u.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expire,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(key))
}
