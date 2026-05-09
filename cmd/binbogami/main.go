package main

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/handlers/auth"
	"github.com/jurgisjaska/binbogami/internal/handlers/public"
	"github.com/jurgisjaska/binbogami/internal/handlers/v1"
	"github.com/jurgisjaska/binbogami/internal/handlers/v1/user"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	config, err := internal.CreateConfig()
	if err != nil {
		log.Fatalln("configuration failure")
	}

	database, err := internal.ConnectDatabase(config.Database)
	if err != nil {
		log.Fatalln("database failure")
	}
	defer func() { _ = database.Close() }()

	// create mail dialer
	dialer := internal.CreateDialer(config.Mail)

	e := echo.New()
	e.Use(middleware.RequestLogger())
	// @todo if this ever goes to production it needs to have proper values!
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	e.HTTPErrorHandler = api.CustomHTTPErrorHandler
	e.Validator = &api.Validator{Validator: validator.New()}
	auth.CreateAuth(e, database, config, dialer)

	// public resources that are not related with auth
	// must be accessible without authentication
	pg := e.Group("/public")
	public.CreatePublic(pg, database)

	// main API
	g := e.Group("/v1")
	g.Use(echojwt.WithConfig(token.CreateJWTConfig(config.Secret)))

	user.CreateUser(g, database)
	user.CreateConfiguration(g, database)
	// v1.CreateInvitation(g, database, config, dialer)

	v1.CreateBook(g, database)
	v1.CreateCategory(g, database)
	v1.CreateLocation(g, database)

	v1.CreateEntry(g, database)

	if err := e.Start(fmt.Sprintf(":%d", config.App.Port)); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
