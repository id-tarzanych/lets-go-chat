package user

import (
	"github.com/id-tarzanych/lets-go-chat/internal/types"
	"github.com/id-tarzanych/lets-go-chat/models"
)

type UserDao interface {
	Create(u *models.User) error
	Update(u *models.User) error
	Delete(id types.Uuid) error
	GetById(id types.Uuid) (models.User, error)
	GetByUserName(username string) (models.User, error)
	GetAll() ([]models.User, error)
}
