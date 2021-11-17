package user

import (
	"context"
	"gorm.io/gorm"

	"github.com/id-tarzanych/lets-go-chat/internal/types"
	"github.com/id-tarzanych/lets-go-chat/models"
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
	db *gorm.DB
}

func NewDatabaseUserRepository(db *gorm.DB) (*DatabaseUserRepository, error) {
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}

	return &DatabaseUserRepository{db}, nil
}

func (d DatabaseUserRepository) Create(ctx context.Context, u *models.User) error {
	if result := d.db.Create(&u); result.Error != nil {
		return result.Error
	}

	return nil
}

func (d DatabaseUserRepository) Update(ctx context.Context, u *models.User) error {
	result := d.db.Model(&u).Updates(models.User{UserName: u.UserName, PasswordHash: u.PasswordHash})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (d DatabaseUserRepository) Delete(ctx context.Context, id types.Uuid) error {
	if result := d.db.Delete(&models.User{}, id); result.Error != nil {
		return result.Error
	}

	return nil
}

func (d DatabaseUserRepository) GetById(ctx context.Context, id types.Uuid) (models.User, error) {
	u := models.User{}

	result := d.db.First(&u, id)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	return u, nil
}

func (d DatabaseUserRepository) GetByUserName(ctx context.Context, name string) (models.User, error) {
	u := models.User{}

	result := d.db.Where("username = ?", name).First(&u)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	return u, nil
}

func (d DatabaseUserRepository) GetAll(context.Context) ([]models.User, error) {
	var users []models.User

	result := d.db.Find(&users)
	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}
