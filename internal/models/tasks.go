package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	ID          int
	Title       string `json:"title"`
	Description string `json:"description"`
	EndDate     string `json:"end_date"`
	UserId      int64  `json:"user_id"`
}

type DBFirstError struct{}

func (m *DBFirstError) Error() string {
	return "Нет таких задач!"
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) Create(task Task) error {
	result := m.Db.Create(&task)

	return result.Error
}

func (m *TaskModel) FindSame(id int, title, description string) (*Task, error) {
	existTask := Task{}

	if id != 0 {
		result := m.Db.First(&existTask, Task{ID: id})
		if result.Error != nil {
			return nil, result.Error
		}
	} else if title != "" {
		result := m.Db.First(&existTask, Task{Title: title, Description: description})
		if result.Error != nil {
			return nil, result.Error
		}
	} else {
		return nil, &DBFirstError{}
	}

	return &existTask, nil
}

func (m *TaskModel) GetAllByUserId(userId int64) ([]Task, error) {
	var userTasks []Task
	result := m.Db.Find(&userTasks, Task{UserId: userId})

	if result.Error != nil {
		return nil, result.Error
	}

	return userTasks, nil
}

func (m *TaskModel) DropTask(task Task) error {
	result := m.Db.Delete(&task)

	return result.Error
}
