package models

import (
	"gorm.io/gorm"
	"time"

	"github.com/id-tarzanych/lets-go-chat/internal/types"
)

type Token struct {
	gorm.Model

	Token      string `gorm:"primaryKey,column:token"`
	UserId     types.Uuid
	Expiration time.Time
}

func NewToken(token string, userId types.Uuid, expiration time.Time) *Token {
	return &Token{Token: token, UserId: userId, Expiration: expiration}
}
