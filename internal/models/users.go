package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name       string `json:"name"`
	TelegramId int64  `json:"telegram_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ChatId     int64  `json:"chat_id"`
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

func (m *UserModel) GetAllTasks() ([]Task, error) {
	tasks := make([]Task, 0)
	err := m.Db.Preload("users").Find(&tasks).Error

	return tasks, err
}
