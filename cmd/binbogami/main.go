package main

import (
	"log"

	"github.com/jurgisjaska/binbogami/app"
	"github.com/jurgisjaska/binbogami/app/api/token"
	"github.com/jurgisjaska/binbogami/app/handler"
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

	v1 := e.Group("/v1")
	v1.Use(echojwt.WithConfig(token.CreateJWTConfig(config.Secret)))

	// temporary disabled before refactor
	// v1.CreateOrganization(e, database)
	// v1.CreateBook(e, database)
	// v1.CreateCategory(e, database)

	e.Logger.Fatal(e.Start(":8001"))
}
