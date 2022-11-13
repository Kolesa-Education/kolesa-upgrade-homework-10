package models

import (
	"gorm.io/gorm"
	"time"
)

type Task struct {
	ID          uint
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EndDate     time.Time `json:"end_date"`
	UserID      int64     `json:"user_Id"`
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) CreateTask(task Task) error {
	result := m.Db.Create(&task)
	return result.Error
}

func (m *TaskModel) DeleteTask(taskId int64, userId int64) error {
	db := m.Db.Where("user_id = ?", userId).Where("id = ?", taskId).Delete(&Task{})
	return db.Error
}
func (m *TaskModel) GetAll(userId int64) ([]Task, error) {
	var tasks []Task
	result := m.Db.Find(&tasks, Task{UserID: userId})
	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}
