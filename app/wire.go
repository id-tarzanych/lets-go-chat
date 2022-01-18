//go:build ignore
// +build ignore

package app

import (
	"os"

	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/id-tarzanych/lets-go-chat/configurations"
	"github.com/id-tarzanych/lets-go-chat/db"
	"github.com/id-tarzanych/lets-go-chat/db/message"
	"github.com/id-tarzanych/lets-go-chat/db/token"
	"github.com/id-tarzanych/lets-go-chat/db/user"
)

func InitializeApp(config *configurations.Configuration) (Application, error) {
	wire.Build(
		ProvideApp,
		ProvideLogger,
		ProvideDb,
		ProvideUserRepo,
		ProvideTokenRepo,
		ProvideMessageRepo,
	)
	return Application{}, nil
}

func ProvideApp(
	cfg *configurations.Configuration,
	dbPool *gorm.DB,

	logger logrus.FieldLogger,

	userRepo user.UserRepository,
	tokenRepo token.TokenRepository,
	messageRepo message.MessageRepository,
) Application {
	return Application{
		config: cfg,
		db:     dbPool,
		logger: logger,

		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		messageRepo: messageRepo,
	}
}

func ProvideLogger() logrus.FieldLogger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	return logger
}

func ProvideDb(config *configurations.Configuration, logger logrus.FieldLogger) (*gorm.DB, error) {
	switch db.DbType(config.Database.Type) {
	case db.Postgres:
		logger.Println("Using PostgreSQL database...")

		return db.NewPostgresSession(config.Database)
	default:
		logger.Println("Unrecognized database type! Fallback to in-memory database!")

		return db.NewInMemorySession()
	}
}

func ProvideUserRepo(db *gorm.DB) (user.UserRepository, error) {
	return user.NewDatabaseUserRepository(db)
}

func ProvideTokenRepo(db *gorm.DB) (token.TokenRepository, error) {
	return token.NewDatabaseTokenRepository(db)
}

func ProvideMessageRepo(db *gorm.DB) (message.MessageRepository, error) {
	return message.NewDatabaseMessageRepository(db)
}
