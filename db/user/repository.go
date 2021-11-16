package user

import (
	"context"
	"database/sql"
	"github.com/id-tarzanych/lets-go-chat/internal/types"
	"github.com/id-tarzanych/lets-go-chat/models"
	"log"
)

type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	Update(ctx context.Context, u *models.User) error
	Delete(ctx context.Context, id types.Uuid) error
	GetById(ctx context.Context, id types.Uuid) (models.User, error)
	GetByUserName(ctx context.Context, name string) (models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
}

type DatabaseUserRepository struct {
	dbPool *sql.DB
}

func NewDatabaseUserRepository(dbPool *sql.DB) *DatabaseUserRepository {
	pool := &DatabaseUserRepository{dbPool}

	return pool
}

func (d DatabaseUserRepository) Create(ctx context.Context, u *models.User) error {
	stmt, err := d.dbPool.Prepare(
		`INSERT INTO users (uuid, username, password) VALUES ($1, $2, $3)`,
	)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(string(u.Id()), u.UserName(), u.PasswordHash())
	if err != nil {
		return err
	}

	log.Println(result)

	return nil
}

func (d DatabaseUserRepository) Update(ctx context.Context, u *models.User) error {
	stmt, _ := d.dbPool.Prepare(`UPDATE users
		SET username = $1, password = $2
		WHERE uuid = $3
	`)

	_, err := stmt.Exec(
		u.UserName(),
		u.PasswordHash(),
		string(u.Id()),
	)
	if err != nil {
		return err
	}

	return nil
}

func (d DatabaseUserRepository) Delete(ctx context.Context, id types.Uuid) error {
	stmt, _ := d.dbPool.Prepare(`DELETE FROM users WHERE uuid = $1`)

	if _, err := stmt.Exec(id); err != nil {
		return err
	}

	return nil
}

func (d DatabaseUserRepository) GetById(ctx context.Context, id types.Uuid) (models.User, error) {
	u := models.User{}
	var userId types.Uuid
	var username, hash string

	err := d.dbPool.QueryRow("SELECT * FROM USERS WHERE uuid = $1", id).Scan(&userId, &username, &hash)
	if err != nil {
		return u, err
	}

	u.SetId(userId).SetUserName(username).SetPasswordHash(hash)

	return u, nil
}

func (d DatabaseUserRepository) GetByUserName(ctx context.Context, name string) (models.User, error) {
	u := models.User{}
	var userId types.Uuid
	var username, hash string

	err := d.dbPool.QueryRow("SELECT * FROM USERS WHERE username = $1", name).Scan(&userId, &username, &hash)
	if err != nil {
		return u, err
	}

	u.SetId(userId).SetUserName(username).SetPasswordHash(hash)

	return u, nil
}

func (d DatabaseUserRepository) GetAll(context.Context) ([]models.User, error) {
	// This now uses the unexported global variable.
	rows, err := d.dbPool.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var userId types.Uuid
		var username, hash string

		err := rows.Scan(&userId, &username, &hash)
		if err != nil {
			return nil, err
		}

		u := models.User{}
		u.SetId(userId).SetUserName(username).SetPasswordHash(hash)

		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
