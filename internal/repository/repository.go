package repository

import "gorm.io/gorm"

type Repository struct {
	User
	Task
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User: newUserRepos(db),
		Task: newTaskRepos(db),
	}
}
