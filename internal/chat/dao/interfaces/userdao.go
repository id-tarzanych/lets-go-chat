package interfaces

import (
	"github.com/id-tarzanych/lets-go-chat/internal/chat/models"
	"github.com/id-tarzanych/lets-go-chat/internal/chat/types"
)

type UserDao interface {
	Create (u *models.User) error
	Update(u *models.User) error
	Delete(id types.Uuid) error
	GetById(id types.Uuid) (models.User, error)
	GetByUserName(username string) (models.User, error)
	GetAll() ([]models.User, error)
}