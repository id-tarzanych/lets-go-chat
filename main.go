package main

import (
	"fmt"
	"os"

	"github.com/id-tarzanych/lets-go-chat/api/server"
	"github.com/id-tarzanych/lets-go-chat/db/user"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	userDao := inmemory.MakeUserDao("inmemory")
	srv := server.New(&userDao)

	srv.Handle()

	return nil
}
