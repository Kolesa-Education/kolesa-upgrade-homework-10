package bot

import (
	"log"
	"time"
	"upgrade/internal/models"

	"gopkg.in/telebot.v3"
)

type TodoBot struct {
	Bot   *telebot.Bot
	Users *models.UserModel
}

//var winSticker = &telebot.Sticker{
//	File: telebot.File{
//		FileID: "CAACAgIAAxkBAAEGMEZjVspD4JulorxoH7nIwco5PGoCsAACJwADr8ZRGpVmnh4Ye-0RKgQ",
//	},
//	Width:    512,
//	Height:   512,
//	Animated: true,
//}
//
//var loseSticker = &telebot.Sticker{
//	File: telebot.File{
//		FileID: "CAACAgIAAxkBAAEGUfFjZnDTpRmHFgABq_zncY60TpKZhlUAAgsBAAIWfGgDjgzxiPsp7OIrBA",
//	},
//	Width:    512,
//	Height:   512,
//	Animated: true,
//}

func (bot *TodoBot) StartHandler(ctx telebot.Context) error {
	newUser := models.User{
		Name:       ctx.Sender().Username,
		TelegramId: ctx.Chat().ID,
		FirstName:  ctx.Sender().FirstName,
		LastName:   ctx.Sender().LastName,
		ChatId:     ctx.Chat().ID,
	}

	existUser, err := bot.Users.FindOne(ctx.Chat().ID)

	if err != nil {
		log.Printf("Ошибка получения пользователя %v", err)
	}

	if existUser == nil {
		err := bot.Users.Create(newUser)

		if err != nil {
			log.Printf("Ошибка создания пользователя %v", err)
		}
	}
	return ctx.Send("Привет, " + ctx.Sender().FirstName)
}

func (bot *TodoBot) CreateTodoHandler(ctx telebot.Context) error {
	taskArgs := ctx.Args()

	if len(taskArgs) == 0 {
		return ctx.Send("Вы не ввели ни одной задачи")
	}

	existUser, err := bot.Users.FindOne(ctx.Chat().ID)
	if err != nil {
		log.Printf("Ошибка получения пользователя %v", err)
	}

	newTask := models.Task{
		Title:       taskArgs[0],
		Description: taskArgs[0],
		EndDate:     "11-10-2022",
		UserID:      existUser.ID,
	}
	return ctx.Send("Создана задача", newTask.Title, newTask.Description)
}

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
