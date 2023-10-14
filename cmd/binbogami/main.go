package main

import (
	"fmt"
	"log"

	"github.com/jurgisjaska/binbogami/app"
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

	fmt.Println("binbogami")
}
