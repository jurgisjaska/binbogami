package main

import (
	"log"

	"github.com/jurgisjaska/binbogami/app"
	"github.com/jurgisjaska/binbogami/app/handler"
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
	defer database.Close()

	e := echo.New()
	handler.CreateCategory(e, database)
	handler.CreateOrganization(e, database)

	e.Logger.Fatal(e.Start(":8001"))
}
