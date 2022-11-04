package models

type User struct {
	name       string
	telegramId int
	firstName  string
	lastName   string
	chatId     int
}

type Task struct {
	title       string
	description string
	endDate     string
}
