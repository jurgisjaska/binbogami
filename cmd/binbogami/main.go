package main

import (
	"log"

	"github.com/jurgisjaska/binbogami/app"
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
	handler.CreateAuth(e, database, config)

	g := e.Group("/v1")
	g.Use(echojwt.WithConfig(token.CreateJWTConfig(config.Secret)))

	// temporary disabled before refactor
	v1.CreateOrganization(g, database)
	// g.CreateBook(e, database)
	// g.CreateCategory(e, database)

	e.Logger.Fatal(e.Start(":8001"))
}
