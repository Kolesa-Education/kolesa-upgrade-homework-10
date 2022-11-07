package models

import (
	"gorm.io/gorm"
)

type Tasks struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	EndDate     string `json:"end_date"`
	UserId      uint   `json:"user_id"`
}

type TasksModel struct {
	Db *gorm.DB
}

func (m *TasksModel) Create(task Tasks) error {
	err := m.Db.Create(&task).Error
	return err
}

func (m *TasksModel) DeleteTask(id uint) error {
	return m.Db.Delete(&Tasks{}, id).Error
}

func (m *TasksModel) GetTasks(userId uint) ([]Tasks, error) {
	var tasks []Tasks
	err := m.Db.Find(&tasks, Tasks{UserId: userId}).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
