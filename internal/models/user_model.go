package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name       string
	TelegramId int64
	FirstName  string
	LastName   string
	ChatId     int64
	Tasks      []Task
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
