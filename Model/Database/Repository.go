package Database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type User struct {
	gorm.Model
	Id         int `gorm:"primaryKey"`
	TelegramId int64
	FirstName  string
	LastName   string
	ChatId     int64
	Tasks      []Task `gorm:"foreignKey:UserId"`
}

type Task struct {
	gorm.Model
	Id          int `gorm:"primaryKey"`
	UserId      int
	Title       string
	Description string
	EndDate     time.Time
}

func NewDatabase(dbAddress, dbName, dbUsername, dbPassword string) *gorm.DB {
	dsn := dbUsername + ":" + dbPassword + "@tcp(" + dbAddress + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatal("Error connecting to database: ", err.Error())
	}
	return db
}

func NewUser(database *gorm.DB, telegramId int64, firstName, lastName string, chatId int64) {
	err := database.Create(&User{
		TelegramId: telegramId,
		FirstName:  firstName,
		LastName:   lastName,
		ChatId:     chatId,
	})
	if err == nil {
		log.Printf("Created new user:\n" +
			"\tTelegram ID: %v\n" +
			"\tFirst Name:")
	}
	log.Println(err.Error)
}

func NewTask(database *gorm.DB, userId int, title, description string, endDate time.Time) {
	database.Create(&Task{
		UserId:      userId,
		Title:       title,
		Description: description,
		EndDate:     endDate,
	})
}

func GetUser(database *gorm.DB, id int) {
	var user User
	database.Preload("Tasks.User_Id").First(&user, "id = ?", id)
	println(user.Id, len(user.Tasks))
}

func GetTask(database *gorm.DB, id int) {
	var task Task
	database.First(&task, "id = ?", id)
	fmt.Printf("%v\n", task)
	println(task.UserId)
}
