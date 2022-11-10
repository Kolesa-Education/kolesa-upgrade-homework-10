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

func (t *TaskModel) Delete(id int, nameColumn string) error {
	result := t.Db.Where(nameColumn+" = ?", id).Delete(&Task{})
	return result.Error
}
