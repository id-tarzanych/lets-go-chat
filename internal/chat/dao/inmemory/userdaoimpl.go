package inmemory

import (
	"errors"
	"github.com/id-tarzanych/lets-go-chat/internal/chat/models"
	"github.com/id-tarzanych/lets-go-chat/internal/chat/types"
)

type UserDaoInMemory struct {
	repository map[types.Uuid]models.User
}

func NewUserDaoInMemory() *UserDaoInMemory {
	dao := &UserDaoInMemory{}
	dao.repository = make(map[types.Uuid]models.User, 0)

	return dao
}

func (dao *UserDaoInMemory) Create(u *models.User) error {
	dao.repository[u.Id()] = *u

	return nil
}

func (dao *UserDaoInMemory) Update(u *models.User) error {
	dao.repository[u.Id()] = *u

	return nil
}

func (dao *UserDaoInMemory) Delete(id types.Uuid) error {
	delete(dao.repository, id)

	return nil
}

func (dao *UserDaoInMemory) GetById(id types.Uuid) (models.User, error) {
	user, ok := dao.repository[id];
	if !ok {
		return models.User{}, errors.New("user not found")
	}

	return user, nil
}

func (dao *UserDaoInMemory) GetByUserName(username string) (models.User, error) {
	for _, user := range dao.repository {
		if user.UserName() == username {
			return user, nil
		}
	}

	return models.User{}, errors.New("user not found")
}

func (dao *UserDaoInMemory) GetAll() ([]models.User, error) {
	all := make([]models.User, 0, len(dao.repository))

	for  _, value := range dao.repository {
		all = append(all, value)
	}

	return all, nil
}
