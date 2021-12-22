package testdb

import (
	"github.com/id-tarzanych/lets-go-chat/internal/types"
	"github.com/id-tarzanych/lets-go-chat/models"
	"gorm.io/gorm"
	"time"
)

func Truncate(db *gorm.DB) error {
	result := db.Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.User{})

	if result.Error != nil {
		return result.Error
	}

	result = db.Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Token{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func SeedUsers(db *gorm.DB) (map[types.Uuid]models.User, error) {
	users := []models.User{
		{
			ID:           "6b2db94c-6fce-4673-a1ce-d24ff6bd4d35",
			UserName:     "user1",
			PasswordHash: "ef797c8118f02dfb649607dd5d3f8c7623048c9c063d532cc95c5ed7a898a64f", // 12345678
		},
		{
			ID:           "95a62e6c-e0e7-46ee-8bc3-6cca62b4cb09",
			UserName:     "user2",
			PasswordHash: "2b4bb6f5e5e12daa7d462fbb3bd6d342106ef6fc091073e6086e9771775a38bd", // Secr3tPa$sw0rD
		},
		{
			ID:           "b3341b87-c561-4142-bd28-f9ecde74822b",
			UserName:     "user3",
			PasswordHash: "f386e47a084d97a2a4efa95fbe8aeee5e13d9c0482575130080d973289a129cc", // S@ntaClau$
		},
	}

	result := db.Create(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	usersMap := make(map[types.Uuid]models.User)
	for i := range users {
		usersMap[users[i].ID] = users[i]
	}

	return usersMap, nil
}

func SeedTokens(db *gorm.DB) (map[string]models.Token, error) {
	_, err := SeedUsers(db)
	if err != nil {
		return nil, err
	}

	tokens := []models.Token{
		{
			Token:      "18sqhpLyANr7ypoK",
			UserId:     "6b2db94c-6fce-4673-a1ce-d24ff6bd4d35",
			Expiration: time.Now().Add(time.Hour * 24),
		},
		{
			Token:      "8n9hKwlT9l037PZb",
			UserId:     "6b2db94c-6fce-4673-a1ce-d24ff6bd4d35",
			Expiration: time.Now().Add(time.Hour * 24),
		},
		{
			Token:      "qGVCS6GujpFP0zyP",
			UserId:     "b3341b87-c561-4142-bd28-f9ecde74822b",
			Expiration: time.Now().Add(time.Hour * 24),
		},
	}

	result := db.Create(&tokens)

	if result.Error != nil {
		return nil, result.Error
	}

	tokensMap := make(map[string]models.Token)
	for i := range tokens {
		tokensMap[tokens[i].Token] = tokens[i]
	}

	return tokensMap, nil
}
