package api

import (
	"github.com/golang-jwt/jwt"
	"github.com/jurgisjaska/binbogami/app"
)

type (
	TokenClaims struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		jwt.StandardClaims
	}

	Token struct {
		configuration *app.Config
	}
)

func (t *Token) CreateToken(c TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, c)

	return token.SignedString(t.configuration.SecretSalt)
}

func (t *Token) CreateRefreshToken(c jwt.StandardClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, c)

	return token.SignedString(t.configuration.SecretSalt)
}

func (t *Token) ParseToken(ts string) *TokenClaims {
	token, _ := jwt.ParseWithClaims(ts, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.configuration.SecretSalt), nil
	})

	return token.Claims.(*TokenClaims)
}

func (t *Token) ParseRefreshToken(rt string) *jwt.StandardClaims {
	token, _ := jwt.ParseWithClaims(rt, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.configuration.SecretSalt), nil
	})

	return token.Claims.(*jwt.StandardClaims)
}

func CreateToken(c *app.Config) *Token {
	return &Token{configuration: c}
}
