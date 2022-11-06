package Database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

type User struct {
	gorm.Model
	id         int   `gorm:"primaryKey"`
	TelegramId int64 `gorm:"column:Telegram_Id"`
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

type Database struct {
	Connection *gorm.DB
	DbAddress  string
	DbName     string
	DbUsername string
	DbPassword string
}

func NewDatabase(dbAddress, dbName, dbUsername, dbPassword string) *Database {
	dsn := dbUsername + ":" + dbPassword + "@tcp(" + dbAddress + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err.Error())
	}
	return &Database{
		Connection: db,
		DbAddress:  dbAddress,
		DbName:     dbName,
		DbUsername: dbUsername,
		DbPassword: dbPassword,
	}
}

func (database *Database) NewUser(telegramId int64, firstName, lastName string, chatId int64) error {
	//Check if user already exists
	dbFind := database.Connection.Find(&User{}, "Telegram_Id = ?", telegramId)
	if dbFind.Error != nil {
		return dbFind.Error
	}
	println(dbFind.RowsAffected)
	if dbFind.RowsAffected > 0 {
		return nil
	}

	dbCreate := database.Connection.Create(&User{
		TelegramId: telegramId,
		FirstName:  firstName,
		LastName:   lastName,
		ChatId:     chatId,
	})
	return dbCreate.Error
}

func (database *Database) NewTask(userId int64, args []string) error {
	//Parse time in string using defined format and timezone
	loc, _ := time.LoadLocation("Local")
	endDate, err := time.ParseInLocation("02.01.2006 15:04", args[2], loc)
	if err != nil {
		return err
	}
	db := database.Connection.Create(&Task{
		UserId:      userId,
		Title:       args[0],
		Description: args[1],
		EndDate:     endDate,
	})
	return db.Error
}

// Get user with slice of his tasks
func (database *Database) GetUser(telegramId int) error { //Не потребовалось, возможность для расширения
	var user User
	database.Connection.Preload("Tasks").First(&user, "Telegram_Id = ?", telegramId)
	return database.Connection.Error
}

func (database *Database) GetUserTasks(userId int64) ([]Task, error) {
	var tasks []Task
	db := database.Connection.Find(&tasks, "User_Id = ?", userId)
	return tasks, db.Error
}

func (database *Database) DeleteTask(taskId int, userId int64) (int64, error) {
	db := database.Connection.Where("User_Id = ?", userId).Delete(&Task{}, taskId)
	return db.RowsAffected, db.Error
}
