package app

import (
	"database/sql"
	"fmt"
	"github.com/id-tarzanych/lets-go-chat/db"
	"github.com/id-tarzanych/lets-go-chat/db/user"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/id-tarzanych/lets-go-chat/configurations"
	_ "github.com/lib/pq"
)

type Application struct {
	config   *configurations.Configuration
	dbPool   *db.AppDBPool
	userRepo *user.DatabaseUserRepository
}

func New(cfg *configurations.Configuration) (*Application, error) {
	appDBPool := initDBPool(cfg)

	app := Application{
		config:   cfg,
		dbPool:   &appDBPool,
		userRepo: user.NewDatabaseUserRepository(appDBPool.GetDB()),
	}

	appDBPool.InitDatabase()

	return &app, nil
}

func initDBPool(cfg *configurations.Configuration) db.AppDBPool {
	var dsn string

	switch d := cfg.Database; d.Type {
	case db.Postgres:
		pool := db.PostgresPool{}
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s", d.User, d.Password, d.Host, d.Port, d.Database)

		parameters := make([]string, 0)
		if d.Ssl {
			parameters = append(parameters, "sslmode=require")
		}

		if len(parameters) > 0 {
			dsn += "?" + strings.Join(parameters, "&")
		}

		sqlPool, err := sql.Open(string(d.Type), dsn)
		if err != nil {
			log.Fatal(err)
		}

		err = sqlPool.Ping()
		if err != nil {
			log.Fatal(err)
		}

		pool.DB = sqlPool

		return &pool
	}

	return nil
}

func (a *Application) Config() *configurations.Configuration {
	return a.config
}

func (a *Application) DBPool() *db.AppDBPool {
	return a.dbPool
}

func (a *Application) UserRepo() *user.DatabaseUserRepository {
	return a.userRepo
}
