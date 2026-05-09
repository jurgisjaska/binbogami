package token

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateToken(t *testing.T) {
	key := "test-secret-key"
	userId := uuid.New()
	userEmail := "test@example.com"

	u := &user.User{
		Id:      userId,
		Email:   userEmail,
		Name:    "John",
		Surname: "Doe",
	}

	tokenStr, err := CreateToken(u, key)
	require.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	// Validate the created token
	parsedToken, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	require.NoError(t, err)

	claims, ok := parsedToken.Claims.(*Claims)
	require.True(t, ok)
	require.NotNil(t, claims.Id)
	require.NotNil(t, claims.Email)

	assert.Equal(t, userId, *claims.Id)
	assert.Equal(t, userEmail, *claims.Email)
	assert.Equal(t, "John Doe", claims.Name)

	// Check expiration (should be roughly 72 hours from now)
	expirationTime := claims.ExpiresAt.Time
	expectedExpiration := time.Now().Add(time.Hour * expire)
	assert.WithinDuration(t, expectedExpiration, expirationTime, time.Minute) // 1 min delta
}

func TestCreateJWTConfig(t *testing.T) {
	key := "test-secret-key"
	config := CreateJWTConfig(key)

	assert.Equal(t, []byte(key), config.SigningKey)
	require.NotNil(t, config.NewClaimsFunc)

	// Test NewClaimsFunc
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	claims := config.NewClaimsFunc(c)
	assert.IsType(t, &Claims{}, claims)
}

func TestFromContext(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	userId := uuid.New()
	userEmail := "test@example.com"
	claims := &Claims{
		Id:    &userId,
		Email: &userEmail,
		Name:  "Jane Doe",
	}

	token := &jwt.Token{
		Claims: claims,
	}

	// Set the token in context as echo-jwt middleware does
	c.Set("user", token)

	retrievedClaims := FromContext(c)

	require.NotNil(t, retrievedClaims)
	assert.Equal(t, userId, *retrievedClaims.Id)
	assert.Equal(t, userEmail, *retrievedClaims.Email)
	assert.Equal(t, "Jane Doe", retrievedClaims.Name)
}
