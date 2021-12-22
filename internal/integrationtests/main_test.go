package integrationtests

import (
	"os"
	"strconv"
	"testing"

	"github.com/id-tarzanych/lets-go-chat/app"
	"github.com/id-tarzanych/lets-go-chat/configurations"
	log "github.com/sirupsen/logrus"
)

var (
	a *app.Application
)

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	port, _ := strconv.Atoi(os.Getenv("LETS_GO_CHAT_DATABASE__PORT"))
	cfg := &configurations.Configuration{
		Database: configurations.Database{
			Type:     "postgres",
			Host:     os.Getenv("LETS_GO_CHAT_DATABASE__HOST"),
			Port:     port,
			Protocol: "",
			User:     os.Getenv("LETS_GO_CHAT_DATABASE__USER"),
			Password: os.Getenv("LETS_GO_CHAT_DATABASE__PASSWORD"),
			Database: os.Getenv("LETS_GO_CHAT_DATABASE__DATABASE"),
		},
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	a = application

	return m.Run()
}
