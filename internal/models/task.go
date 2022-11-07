package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	UserID      uint
	Title       string
	Description string
	EndDate     time.Time
}
