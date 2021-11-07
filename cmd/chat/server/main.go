package main

import (
	"fmt"
	"os"

	"github.com/id-tarzanych/lets-go-chat/internal/chat/dao/factory"
	"github.com/id-tarzanych/lets-go-chat/internal/chat/server"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	userDao := factory.MakeUserDao("inmemory")
	srv := server.New(&userDao)

	srv.Handle()

	return nil
}
