package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	ID          uint   `json:"id,omitempty" gorm:"primaryKey:id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	EndDate     string `json:"end_date"`
	UserId      uint   `json:"user_id"`
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) Create(tasks Task) error {
	result := m.Db.Create(&tasks)

	return result.Error
}

func (m *TaskModel) FindOne(id uint) (*Task, error) {
	t := Task{}
	result := m.Db.First(&t, Task{ID: id})
	if result.Error != nil {
		return nil, result.Error
	}
	return &t, nil
}
