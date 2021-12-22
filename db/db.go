package db

import (
	"fmt"

	"github.com/id-tarzanych/lets-go-chat/configurations"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DbType string

const (
	Postgres DbType = "postgres"
	SQLite   DbType = "sqlite"
)

func NewPostgresSession(cfg configurations.Database) (*gorm.DB, error) {
	var sslMode string

	switch cfg.Ssl {
	case true:
		sslMode = "require"
	default:
		sslMode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", cfg.Host, cfg.User, cfg.Password, cfg.Database, cfg.Port, sslMode)
	gormPool, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return gormPool, nil
}

func NewInMemorySession() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
}
