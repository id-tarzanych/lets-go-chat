package app

import (
	"github.com/id-tarzanych/lets-go-chat/db/token"
	"github.com/sirupsen/logrus"
	"log"
	"os"

	"gorm.io/gorm"

	"github.com/id-tarzanych/lets-go-chat/configurations"
	"github.com/id-tarzanych/lets-go-chat/db"
	"github.com/id-tarzanych/lets-go-chat/db/user"
)

type Application struct {
	config *configurations.Configuration
	db     *gorm.DB
	logger *logrus.Logger

	userRepo  *user.DatabaseUserRepository
	tokenRepo *token.DatabaseTokenRepository
}

func New(cfg *configurations.Configuration) (*Application, error) {
	dbPool, err := initDB(cfg)

	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	if err != nil {
		logger.Fatal(err)
	}

	userRepo, err := user.NewDatabaseUserRepository(dbPool)
	if err != nil {
		logger.Fatal(err)
	}

	tokenRepo, err := token.NewDatabaseTokenRepository(dbPool)
	if err != nil {
		logger.Fatal(err)
	}

	app := Application{
		config: cfg,
		db:     dbPool,
		logger: logger,

		userRepo:  userRepo,
		tokenRepo: tokenRepo,
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

func (a *Application) Logger() *logrus.Logger  {
	return a.logger
}

func (a *Application) UserRepo() *user.DatabaseUserRepository {
	return a.userRepo
}

func (a *Application) TokenRepo() *token.DatabaseTokenRepository {
	return a.tokenRepo
}
