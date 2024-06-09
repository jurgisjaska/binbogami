package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/handler/auth"
	"github.com/jurgisjaska/binbogami/internal/handler/p"
	"github.com/jurgisjaska/binbogami/internal/handler/v1"
	"github.com/jurgisjaska/binbogami/internal/handler/v1/user"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	// @deprecated use dialer and mail services
	mail, err := internal.ConnectMail(config.Mail)
	if err != nil {
		log.Fatalln("mail failure")
	}
	defer func() { _ = mail.Close() }()

	// create mail dialer
	dialer := internal.CreateDialer(config.Mail)

	e := echo.New()
	// @todo if this ever goes to production it needs to have proper values!
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	e.HTTPErrorHandler = customHTTPErrorHandler // @todo move to the api?
	e.Validator = &api.Validator{Validator: validator.New()}
	auth.CreateAuth(e, database, config, dialer)

	// public resources that are not related with auth
	// mus be accessible without authentication
	pg := e.Group("/p")
	p.CreateInvitation(pg, database)
	p.CreateReset(pg, database)

	// main API
	g := e.Group("/v1")
	g.Use(echojwt.WithConfig(token.CreateJWTConfig(config.Secret)))

	v1.CreateOrganization(g, database)
	user.CreateUser(g, database)
	user.CreateConfiguration(g, database)
	v1.CreateInvitation(g, database, mail, config)
	v1.CreateMember(g, database)

	v1.CreateBook(g, database)
	v1.CreateCategory(g, database)
	v1.CreateLocation(g, database)

	v1.CreateEntry(g, database)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.App.Port)))
}

// customHTTPErrorHandler handles HTTP errors and provides custom error responses.
func customHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	he, ok := err.(*echo.HTTPError)
	if ok {
		if he.Internal != nil {
			if herr, ok := he.Internal.(*echo.HTTPError); ok {
				he = herr
			}
		}
	} else {
		he = &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}
	}

	code := he.Code
	message := he.Message.(string)

	if err == sql.ErrNoRows {
		code = http.StatusNotFound
		message = http.StatusText(code)
	}

	if c.Request().Method == http.MethodHead {
		err = c.NoContent(he.Code)
	} else {
		err = c.JSON(code, api.Error(message))
	}

	if err != nil {
		log.Fatal(err)
	}
}
