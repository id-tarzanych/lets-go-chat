package main

import (
	"fmt"
	"log"
	"os"

	"github.com/id-tarzanych/lets-go-chat/api/server"
	"github.com/id-tarzanych/lets-go-chat/app"
	"github.com/id-tarzanych/lets-go-chat/configurations"
)

func main() {
	config, err := configurations.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	application, err := app.New(config)
	if err != nil {
		log.Fatal(err)
	}

	runServer(application)
}

func runServer(app *app.Application) {
	srv := server.New(*app.Config(), app.UserRepo())
	srv.Handle()
}
