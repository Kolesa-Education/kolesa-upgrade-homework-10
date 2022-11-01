package models

import (
	"database/sql"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title       string `json:"Title"`
	Description string `json:"DESCRIPTION"`
	UserId      int64  `json:"USER_ID"`
	EndDate     int64  `json:"END_DATE"`
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) Create(task Task) error {
	result := m.Db.Create(&task)
	return result.Error
}
func (m *TaskModel) GetAll(taskid int64) *sql.Rows {

	rows, _ := m.Db.Model(&Task{}).Select("tasks.id, tasks.title, tasks.DESCRIPTION, tasks.END_DATE").Joins("join users on users.telegram_id = tasks.USER_ID and tasks.USER_ID=?", taskid).Rows()

	return rows
}
func (m *TaskModel) DeleteTask(taskid string, userid int64) bool {
	m.Db.Where("USER_ID = ?", userid).Where("id = ?", taskid).Delete(&Task{})

	return true
}
