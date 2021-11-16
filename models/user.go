package models

import (
	"github.com/google/uuid"
	"github.com/id-tarzanych/lets-go-chat/internal/types"
	"github.com/id-tarzanych/lets-go-chat/pkg/hasher"
)

type User struct {
	id           types.Uuid
	userName     string
	passwordHash string
}

func NewUser(username, password string) *User {
	id, _ := uuid.NewUUID()
	u := &User{id: types.Uuid(id.String()), userName: username}

	u.SetPassword(password)

	return u
}

func (u User) Id() types.Uuid {
	return u.id
}

func (u User) UserName() string {
	return u.userName
}

func (u User) PasswordHash() string {
	return u.passwordHash
}

func (u *User) SetId(id types.Uuid) *User {
	u.id = id

	return u
}

func (u *User) SetUserName(username string) *User {
	u.userName = username

	return u
}

func (u *User) SetPassword(password string) *User {
	hash, _ := hasher.HashPassword(password)
	u.passwordHash = hash

	return u
}

func (u *User) SetPasswordHash(hash string) *User {
	u.passwordHash = hash

	return u
}
