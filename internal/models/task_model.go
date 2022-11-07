package models

import "gorm.io/gorm"

type Task struct {
	Id          uint
	Title       string
	Description string
	EndDate     string
	UserId      uint
}

type TaskModel struct {
	Db *gorm.DB
}

func (t *TaskModel) Create(task Task) error {
	result := t.Db.Create(&task)

	return result.Error
}
