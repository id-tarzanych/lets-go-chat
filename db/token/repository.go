package token

import (
	"context"
	"gorm.io/gorm"

	"github.com/id-tarzanych/lets-go-chat/internal/types"
	"github.com/id-tarzanych/lets-go-chat/models"
)

type TokenRepository interface {
	Create(ctx context.Context, u *models.Token) error
	Delete(ctx context.Context, token string) error
	Get(ctx context.Context, token string) (models.Token, error)
	GetByUserId(ctx context.Context, userId types.Uuid) ([]models.Token, error)
	GetAll(ctx context.Context) ([]models.Token, error)
}

type DatabaseTokenRepository struct {
	db *gorm.DB
}

func NewDatabaseTokenRepository(db *gorm.DB) (*DatabaseTokenRepository, error) {
	err := db.AutoMigrate(&models.Token{})
	if err != nil {
		return nil, err
	}

	return &DatabaseTokenRepository{db}, nil
}

func (d DatabaseTokenRepository) Create(ctx context.Context, t *models.Token) error {
	if result := d.db.Create(&t); result.Error != nil {
		return result.Error
	}

	return nil
}

func (d DatabaseTokenRepository) Delete(ctx context.Context, token string) error {
	if result := d.db.Delete(&models.Token{}, "token = ?", token); result.Error != nil {
		return result.Error
	}

	return nil
}

func (d DatabaseTokenRepository) Get(ctx context.Context, token string) (models.Token, error) {
	t := models.Token{}

	result := d.db.First(&t, "token = ?", token)
	if result.Error != nil {
		return models.Token{}, result.Error
	}

	return t, nil
}

func (d DatabaseTokenRepository) GetByUserId(ctx context.Context, userId types.Uuid) ([]models.Token, error) {
	var tokens []models.Token

	result := d.db.Where("user_id = ?", userId).Find(&tokens)
	if result.Error != nil {
		return tokens, result.Error
	}

	return tokens, nil
}

func (d DatabaseTokenRepository) GetAll(context.Context) ([]models.Token, error) {
	var tokens []models.Token

	result := d.db.Find(&tokens)
	if result.Error != nil {
		return tokens, result.Error
	}

	return tokens, nil
}
