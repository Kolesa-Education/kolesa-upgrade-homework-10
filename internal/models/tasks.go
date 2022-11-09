package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	UserID      uint   `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	EndDate     string `json:"deleted_at"`
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) Create(task Task) error {

	result := m.Db.Create(&task)

	return result.Error
}

func (m *TaskModel) GetAll(userID uint) ([]Task, error) {

	list := make([]Task, 0)

	result := m.Db.Find(&list, Task{UserID: userID})

	if result.Error != nil {
		return nil, result.Error
	}

	return list, nil
}

func (m *TaskModel) DeleteById(taskID int) error {

	existTask := Task{}

	result := m.Db.Delete(&existTask, taskID)

	return result.Error
}
