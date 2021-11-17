package app

import (
	"log"

	"gorm.io/gorm"

	"github.com/id-tarzanych/lets-go-chat/configurations"
	"github.com/id-tarzanych/lets-go-chat/db"
	"github.com/id-tarzanych/lets-go-chat/db/user"
)

type Application struct {
	config   *configurations.Configuration
	db       *gorm.DB
	userRepo *user.DatabaseUserRepository
}

func New(cfg *configurations.Configuration) (*Application, error) {
	dbPool, err := initDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	repo, err := user.NewDatabaseUserRepository(dbPool)
	if err != nil {
		log.Fatal(err)
	}

	app := Application{
		config:   cfg,
		db:       dbPool,
		userRepo: repo,
	}

	return &app, nil
}

func initDB(cfg *configurations.Configuration) (*gorm.DB, error) {
	switch db.DbType(cfg.Database.Type) {
	case db.Postgres:
		return db.NewPostgresSession(cfg.Database)
	default:
		log.Println("Unrecognized database type! Fallback to in-memory database!")
		return db.NewInMemorySession()
	}
}

func (a *Application) Config() *configurations.Configuration {
	return a.config
}

func (a *Application) DB() *gorm.DB {
	return a.db
}

func (a *Application) UserRepo() *user.DatabaseUserRepository {
	return a.userRepo
}
