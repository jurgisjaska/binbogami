package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/jurgisjaska/binbogami/app"
	"github.com/jurgisjaska/binbogami/app/api"
	"github.com/jurgisjaska/binbogami/app/api/token"
	"github.com/jurgisjaska/binbogami/app/handler"
	"github.com/jurgisjaska/binbogami/app/handler/v1"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func main() {
	config, err := app.CreateConfig()
	if err != nil {
		log.Printf("%+e", err)
		log.Fatalln("configuration failure")
	}

	database, err := app.ConnectDatabase(config.Database)
	if err != nil {
		log.Fatalln("database failure")
	}
	defer func() {
		_ = database.Close()
	}()

	e := echo.New()
	e.HTTPErrorHandler = customHTTPErrorHandler
	handler.CreateAuth(e, database, config)

	g := e.Group("/v1")
	g.Use(echojwt.WithConfig(token.CreateJWTConfig(config.Secret)))

	// temporary disabled before refactor
	v1.CreateOrganization(g, database)
	v1.CreateBook(g, database)
	v1.CreateCategory(g, database)

	e.Logger.Fatal(e.Start(":8001"))
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
