package models

import (
	"gorm.io/gorm"

	"github.com/id-tarzanych/lets-go-chat/internal/types"
)

type Message struct {
	gorm.Model

	AuthorUuid types.Uuid
	Author     User `gorm:"foreignKey:AuthorUuid"`
	Message    string
}
