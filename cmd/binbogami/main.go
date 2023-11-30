package main

import (
	"log"

	"github.com/jurgisjaska/binbogami/app"
	"github.com/jurgisjaska/binbogami/app/api"
	"github.com/jurgisjaska/binbogami/app/handler"
	"github.com/jurgisjaska/binbogami/app/handler/v1"
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

	_ = api.CreateToken(config)

	e := echo.New()
	handler.CreateAuth(e, database, config)

	// v1 := e.Group("/v1")

	v1.CreateOrganization(e, database)
	v1.CreateBook(e, database)
	v1.CreateCategory(e, database)

	e.Logger.Fatal(e.Start(":8001"))
}
