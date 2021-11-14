package db

import (
	"database/sql"
	"log"
)

type AppDBPool interface {
	GetDB() *sql.DB
	InitDatabase()
}

type Pool struct {
	DB *sql.DB
}

func (p *Pool) InitDatabase() {
	panic("Undefined database type")
}

type PostgresPool struct {
	Pool
}

func (p *Pool) GetDB() *sql.DB  {
	return p.DB
}

func (p *PostgresPool) InitDatabase() {
	_, err := p.DB.Exec(`create table if not exists users
		(
			uuid varchar primary key,
			username varchar,
			password varchar
		);
		
		create unique index if not exists users_uuid_uindex
			on users (uuid);
		
		create unique index  if not exists users_username_uindex
			on users (username);
		`)

	if err != nil {
		log.Fatal(err)
	}
}