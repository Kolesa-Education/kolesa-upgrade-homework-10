package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name       string `json:"name,omitempty"`
	TelegramId int64  `json:"telegram_id,omitempty"`
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	ChatId     int64  `json:"chat_id,omitempty"`
}

type UserModel struct {
	Db *gorm.DB
}

func (m *UserModel) Create(user User) error {
	result := m.Db.Create(&user)
	return result.Error
}

func (m *UserModel) FindOne(telegramId int64) (*User, error) {
	existUser := User{}
	result := m.Db.First(&existUser, User{TelegramId: telegramId})
	if result.Error != nil {
		return nil, result.Error
	}
	return &existUser, nil
}
