package main

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/handlers/auth"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// Auth service provides authentication and authorization functionality,
// including user login, registration, and JWT token management.

func main() {
	log.Println("starting auth service")

	config, err := internal.CreateConfig()
	if err != nil {
		log.Fatalln("configuration failure")
	}

	database, err := internal.ConnectDatabase(config.Database)
	if err != nil {
		log.Fatalln("database failure")
	}
	defer func() { _ = database.Close() }()

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.HTTPErrorHandler = api.CustomHTTPErrorHandler
	e.Validator = &api.Validator{Validator: validator.New()}

	// @todo if this ever goes to production it needs to have proper values!
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	dialer := internal.CreateDialer(config.Mail)
	auth.CreateAuth(e, database, config, dialer)

	if err := e.Start(fmt.Sprintf(":%d", config.Auth.Port)); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
