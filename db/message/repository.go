package message

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/id-tarzanych/lets-go-chat/models"
)

type MessageRepository interface {
	Create(ctx context.Context, u *models.Message) error
	Update(ctx context.Context, u *models.Message) error
	Delete(ctx context.Context, id uint) error
	GetAll(ctx context.Context) ([]models.Message, error)
	GetNewerThan(ctx context.Context, time time.Time) ([]models.Message, error)
}

type DatabaseMessageRepository struct {
	db *gorm.DB
}

func NewDatabaseMessageRepository(db *gorm.DB) (*DatabaseMessageRepository, error) {
	err := db.AutoMigrate(&models.Message{})
	if err != nil {
		return nil, err
	}

	return &DatabaseMessageRepository{db}, nil
}

func (d DatabaseMessageRepository) Create(ctx context.Context, u *models.Message) error {
	if result := d.db.Create(&u); result.Error != nil {
		return result.Error
	}

	return nil
}

func (d DatabaseMessageRepository) Update(ctx context.Context, u *models.Message) error {
	result := d.db.Model(&u).Updates(models.Message{Author: u.Author, Message: u.Message})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (d DatabaseMessageRepository) Delete(ctx context.Context, id uint) error {
	if result := d.db.Delete(&models.Message{}, id); result.Error != nil {
		return result.Error
	}

	return nil
}

func (d DatabaseMessageRepository) GetAll(ctx context.Context) ([]models.Message, error) {
	var messages []models.Message

	result := d.db.Order("created_at").Preload("Author").Find(&messages)
	if result.Error != nil {
		return messages, result.Error
	}

	return messages, nil
}

func (d DatabaseMessageRepository) GetNewerThan(ctx context.Context, time time.Time) ([]models.Message, error) {
	var messages []models.Message

	result := d.db.Where("created_at > ?", time).Preload("Author").Order("created_at").Find(&messages)
	if result.Error != nil {
		return messages, result.Error
	}

	return messages, nil
}
