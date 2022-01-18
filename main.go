package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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

	application, err := app.InitializeApp(config)
	if err != nil {
		log.Fatal(err)
	}

	runServer(&application)
}

func runServer(app *app.Application) {
	s := server.New(*app.Config(), app.UserRepo(), app.TokenRepo(), app.MessageRepo(), app.Logger())
	h := server.Handler(s)

	err := http.ListenAndServe(":"+strconv.Itoa(s.Port()), h)
	if err != nil {
		app.Logger().Fatal("Could not start server")
	}
}
