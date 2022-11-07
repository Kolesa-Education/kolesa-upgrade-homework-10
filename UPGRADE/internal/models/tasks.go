package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Task        string `json:"task"`
	Description string `json:"description"`
	End_date    int    `json:"end_date"`
	TelegramId  int64  `json:"telegram_id"`
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) Create(task Task) error {

	result := m.Db.Create(&task)

	return result.Error
}
func (m TaskModel) AllTask(telegram_id int64) ([]Task, error) {
	var tasks []Task
	db := m.Db.Find(&tasks, "telegram_id = ?", telegram_id)
	return tasks, db.Error
}

// func (m *TaskModel) DeleteTask(telegram_id int64) bool {
// 	return true
// }
