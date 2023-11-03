package api

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/jurgisjaska/binbogami/app"
)

type (
	TokenClaims struct {
		Id    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
		jwt.RegisteredClaims
	}

	Token struct {
		configuration *app.Config
	}
)

func (t *Token) CreateToken(c TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, c)

	return token.SignedString(t.configuration.Salt)
}

func (t *Token) CreateRefreshToken(c jwt.RegisteredClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, c)

	return token.SignedString(t.configuration.Salt)
}

func (t *Token) ParseToken(ts string) *TokenClaims {
	token, _ := jwt.ParseWithClaims(ts, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.configuration.Salt), nil
	})

	return token.Claims.(*TokenClaims)
}

func (t *Token) ParseRefreshToken(rt string) *jwt.RegisteredClaims {
	token, _ := jwt.ParseWithClaims(rt, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.configuration.Salt), nil
	})

	return token.Claims.(*jwt.RegisteredClaims)
}

func CreateToken(c *app.Config) *Token {
	return &Token{configuration: c}
}
