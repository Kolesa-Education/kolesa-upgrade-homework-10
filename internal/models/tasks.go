package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	ID 			uint 		`json:"id"`
	Title      	string 		`json:"title"`
	Description string 		`json:"description"`
	EndDate  	string 		`json:"end_date"`
	UserId      uint   		`json:"user_id"`
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

	return result.Error
}

