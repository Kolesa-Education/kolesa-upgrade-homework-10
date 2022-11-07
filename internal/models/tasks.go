package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	EndDate     string `json:"enddate"`
	TelegramID  int64  `json:"telegramID"`
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) Create(task Task) error {

	result := m.Db.Create(&task)

	return result.Error
}

func (m *TaskModel) Delete(id int) error {
	task := Task{}

	result := m.Db.Delete(&task, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (m *TaskModel) GetAll(telegramId int64) ([]Task, error) {
	tasks := []Task{}
	result := m.Db.Where("telegram_id = ?", telegramId).Find(&tasks)

	if result.Error != nil {
		return nil, result.Error
	}

	return tasks, nil
}
