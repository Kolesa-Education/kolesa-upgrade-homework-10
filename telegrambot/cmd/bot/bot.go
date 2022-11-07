package bot

import (
	"log"
	"telegrambot/internal/models"
	"time"

	"gopkg.in/telebot.v3"
)

type UpgradeBot struct {
	Bot   *telebot.Bot
	Users *models.UserModel
	Tasks *models.TaskModel
}

var gameItems = [3]string{
	"камень",
	"ножницы",
	"бумага",
}

// var winSticker = &telebot.Sticker{
// 	File: telebot.File{
// 		FileID: "CAACAgIAAxkBAAEGMEZjVspD4JulorxoH7nIwco5PGoCsAACJwADr8ZRGpVmnh4Ye-0RKgQ",
// 	},
// 	Width:    512,
// 	Height:   512,
// 	Animated: true,
// }

// var loseSticker = &telebot.Sticker{
// 	File: telebot.File{
// 		FileID: "CAACAgIAAxkBAAEGMEhjVsqoRriJRO_d-hrqguHNlLyLvQACogADFkJrCuweM-Hw5ackKgQ",
// 	},
// 	Width:    512,
// 	Height:   512,
// 	Animated: true,
// }

func (bot *UpgradeBot) StartHandler(ctx telebot.Context) error {
	newUser := models.User{
		Name:       ctx.Sender().Username,
		TelegramId: ctx.Chat().ID,
		FirstName:  ctx.Sender().FirstName,
		LastName:   ctx.Sender().LastName,
		ChatId:     ctx.Chat().ID,
	}

	// existUser, err := bot.Users.FindOne(ctx.Chat().ID)

	// if err != nil {
	// 	log.Printf("Ошибка получения пользователя %v", err)
	// }

	//if existUser == nil {
	err := bot.Users.Create(newUser)

	if err != nil {
		log.Printf("Ошибка создания пользователя %v", err)
	}
	//}

	return ctx.Send("Привет, " + ctx.Sender().FirstName)
}

// func (bot *UpgradeBot) GameHandler(ctx telebot.Context) error {
// 	return ctx.Send("Сыграем в камень-ножницы-бумага " +
// 		"Введи твой вариант в формате /try камень")
// }

// func (bot *UpgradeBot) TryHandler(ctx telebot.Context) error {
// 	attempts := ctx.Args()

// 	if len(attempts) == 0 {
// 		return ctx.Send("Вы не ввели ваш вариант")
// 	}

// 	if len(attempts) > 1 {
// 		return ctx.Send("Вы ввели больше одного варианта")
// 	}

// 	try := strings.ToLower(attempts[0])
// 	botTry := gameItems[rand.Intn(len(gameItems))]

// 	if botTry == "камень" {
// 		switch try {
// 		case "ножницы":
// 			ctx.Send(loseSticker)
// 			ctx.Send("🪨")
// 			return ctx.Send("Камень! Ты проиграл!")
// 		case "бумага":
// 			ctx.Send(winSticker)
// 			ctx.Send("🪨")
// 			return ctx.Send("Камень! Ты выиграл!")
// 		}
// 	}

// 	if botTry == "ножницы" {
// 		switch try {
// 		case "камень":
// 			ctx.Send(winSticker)
// 			ctx.Send("✂️")
// 			return ctx.Send("Ножницы! Ты выиграл!")
// 		case "бумага":
// 			ctx.Send(loseSticker)
// 			ctx.Send("✂️")
// 			return ctx.Send("Ножницы! Ты проиграл!")
// 		}
// 	}

// 	if botTry == "бумага" {
// 		switch try {
// 		case "ножницы":
// 			ctx.Send(winSticker)
// 			ctx.Send("📃")
// 			return ctx.Send("Бумага! Ты выиграл!")
// 		case "камень":
// 			ctx.Send(loseSticker)
// 			ctx.Send("📃")
// 			return ctx.Send("Бумага! Ты проиграл!")
// 		}
// 	}

// 	if botTry == try {
// 		return ctx.Send("Ничья!")
// 	}

// 	return ctx.Send("Кажется вы ввели неверный вариант!")
// }

func InitBot(token string) *telebot.Bot {
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)

	if err != nil {
		log.Fatalf("Ошибка при инициализации бота %v", err)
	}

	return b
}

func (bot *UpgradeBot) AddTaskHandler(ctx telebot.Context) error {
	return ctx.Send("Please, add your task by sending it. " +
		"Write it in a format '/add Title: title_of_task Description: description_of_task EndDate: deadline_of_task '")
}

func (bot *UpgradeBot) AddHandler(ctx telebot.Context) error {
	task := ctx.Args()
	log.Printf(task[0])
	if len(task) == 0 {
		return ctx.Send("You did not write your task")
	}
	if len(task) < 3 {
		return ctx.Send("Task dose not have sufficient information")

	}
	//botTask := "OKKKK" //gameItems[rand.Intn(len(gameItems))]

	newTask := models.Task{
		TelegramUserId: ctx.Chat().ID,
		Title:          task[1],
		Description:    task[3],
		EndDate:        task[5],
	}
	// CREATE -- INTO VALUES()
	err := bot.Tasks.Create(newTask)

	if err != nil {
		log.Printf("Ошибка создания пользователя %v", err)
	}
	//}
	return ctx.Send(task[1] + "Your Task Successfully stored")
}

func (bot *UpgradeBot) TasksHandler(ctx telebot.Context) error {
	//Get User id
	t := ctx.Chat().ID

	//Query databse for information of tasks for a particular user
	//SELECT DISTINCT title description FROM table_list WHERE telegram_user_id =

	data, _ := bot.Tasks.GetTask(t)

	str := ""
	str = str + "Here is the list of your tasks: \n"

	for _, d := range data {
		str = str + " Title: " + d.Title
		str = str + " Description: " + d.Description
		str = str + "\n"

	}

	return ctx.Send(str)
}

func (bot *UpgradeBot) DeleteTaskHandler(ctx telebot.Context) error {
	//DELETE FROM table WHERE telegram_user_id = ;
	//DeleteUserTask
	task := ctx.Args()
	delTask := models.Task{
		TelegramUserId: ctx.Chat().ID,
		Title:          task[1],
		Description:    task[3],
	}
	err := bot.Tasks.DeleteUserTask(delTask)
	if err != nil {
		log.Printf("Ошибка создания пользователя %v", err)
	}
	return ctx.Send("The Task is successfully deleted")
}
