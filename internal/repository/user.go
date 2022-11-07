package repository

import (
	"bot-tasker/internal/models"

	"gorm.io/gorm"
)

type User interface {
	AddUser(user models.User) error
	GetUserByChatId(chatId int64) (*models.User, error)
}

type UserRepos struct {
	db *gorm.DB
}

func newUserRepos(db *gorm.DB) *UserRepos {
	return &UserRepos{
		db: db,
	}
}

func (r *UserRepos) AddUser(user models.User) error {
	res := r.db.Create(&user)
	return res.Error
}

func (r *UserRepos) GetUserByChatId(chatId int64) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, models.User{ChatId: chatId})
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
