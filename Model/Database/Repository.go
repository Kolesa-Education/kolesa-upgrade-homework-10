package Database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

type User struct {
	gorm.Model
	id         int
	TelegramId int64 `gorm:"primaryKey;column:Telegram_Id"`
	FirstName  string
	LastName   string
	ChatId     int64
	Tasks      []Task `gorm:"foreignKey:UserId;references:TelegramId""` //Не потребовалось, возможность для расширения
}

type Task struct {
	gorm.Model
	Id          int   `gorm:"primaryKey"`
	UserId      int64 `gorm:"column:User_Id"`
	Title       string
	Description string
	EndDate     time.Time `gorm:"column:End_Date"`
}

func NewDatabase(dbAddress, dbName, dbUsername, dbPassword string) *gorm.DB {
	dsn := dbUsername + ":" + dbPassword + "@tcp(" + dbAddress + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err.Error())
	}
	return db
}

func NewUser(database *gorm.DB, telegramId int64, firstName, lastName string, chatId int64) error {
	db := database.Create(&User{
		TelegramId: telegramId,
		FirstName:  firstName,
		LastName:   lastName,
		ChatId:     chatId,
	})
	return db.Error
}

func NewTask(database *gorm.DB, userId int64, args []string) error {
	loc, _ := time.LoadLocation("Local")
	endDate, err := time.ParseInLocation("02.01.2006 15:04", args[2], loc)
	if err != nil {
		return err
	}
	db := database.Create(&Task{
		UserId:      userId,
		Title:       args[0],
		Description: args[1],
		EndDate:     endDate,
	})
	return db.Error
}

func GetUser(database *gorm.DB, telegramId int) error { //Не потребовалось, возможность для расширения
	var user User
	database.Preload("Tasks").First(&user, "Telegram_Id = ?", telegramId)
	return database.Error
}

func GetUserTasks(database *gorm.DB, userId int64) ([]Task, error) {
	var tasks []Task
	db := database.Find(&tasks, "User_Id = ?", userId)
	return tasks, db.Error
}

func DeleteTask(database *gorm.DB, taskId int, userId int64) (int64, error) {
	db := database.Where("User_Id = ?", userId).Delete(&Task{}, taskId)
	return db.RowsAffected, db.Error
}
