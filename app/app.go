package app

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/id-tarzanych/lets-go-chat/db/message"
	"github.com/id-tarzanych/lets-go-chat/db/token"

	"gorm.io/gorm"

	"github.com/id-tarzanych/lets-go-chat/configurations"
	"github.com/id-tarzanych/lets-go-chat/db"
	"github.com/id-tarzanych/lets-go-chat/db/user"
)

type Application struct {
	config *configurations.Configuration
	db     *gorm.DB
	logger logrus.FieldLogger

	userRepo    user.UserRepository
	tokenRepo   token.TokenRepository
	messageRepo message.MessageRepository
}

func New(cfg *configurations.Configuration) (*Application, error) {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	dbPool, err := initDB(cfg, logger)

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

	messageRepo, err := message.NewDatabaseMessageRepository(dbPool)
	if err != nil {
		logger.Fatal(err)
	}

	app := Application{
		config: cfg,
		db:     dbPool,
		logger: logger,

		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		messageRepo: messageRepo,
	}

	return &app, nil
}

func initDB(cfg *configurations.Configuration, logger logrus.FieldLogger) (*gorm.DB, error) {
	switch db.DbType(cfg.Database.Type) {
	case db.Postgres:
		logger.Println("Using PostgreSQL database...")

		return db.NewPostgresSession(cfg.Database)
	default:
		logger.Println("Unrecognized database type! Fallback to in-memory database!")

		return db.NewInMemorySession()
	}
}

func (a *Application) Config() *configurations.Configuration {
	return a.config
}

func (a *Application) DB() *gorm.DB {
	return a.db
}

func (a *Application) Logger() logrus.FieldLogger {
	return a.logger
}

func (a *Application) UserRepo() user.UserRepository {
	return a.userRepo
}

func (a *Application) TokenRepo() token.TokenRepository {
	return a.tokenRepo
}

func (a *Application) MessageRepo() message.MessageRepository {
	return a.messageRepo
}
