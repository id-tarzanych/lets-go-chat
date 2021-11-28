package models

import (
	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/id-tarzanych/lets-go-chat/internal/types"
	"github.com/id-tarzanych/lets-go-chat/pkg/hasher"
)

type User struct {
	gorm.Model

	ID           types.Uuid `gorm:"primaryKey"`
	UserName     string     `gorm:"column:username"`
	PasswordHash string     `gorm:"column:password"`
}

func NewUser(username, password string) *User {
	id, _ := uuid.NewUUID()
	u := &User{ID: types.Uuid(id.String()), UserName: username}

	u.SetPassword(password)

	return u
}

func (u *User) SetUserName(username string) *User {
	u.UserName = username

	return u
}

func (u *User) SetPassword(password string) *User {
	hash, _ := hasher.HashPassword(password)
	u.PasswordHash = hash

	return u
}

func (u *User) SetPasswordHash(hash string) *User {
	u.PasswordHash = hash

	return u
}
