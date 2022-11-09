package bot

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"kolesa-upgrade-homework-10/internal/models"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
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

var winSticker = &telebot.Sticker{
	File:     telebot.File{FileID: "CAACAgIAAxkBAAEGMEZjVspD4JulorxoH7nIwco5PGoCsAACJwADr8ZRGpVmnh4Ye-0RKgQ"},
	Width:    512,
	Height:   512,
	Animated: true,
}

var loseSticker = &telebot.Sticker{
	File: telebot.File{
		FileID: "CAACAgIAAxkBAAEGMEhjVsqoRriJRO_d-hrqguHNlLyLvQACogADFkJrCuweM-Hw5ackKgQ",
	},
	Width:    512,
	Height:   512,
	Animated: true,
}

func (bot *UpgradeBot) StartHandler(ctx telebot.Context) error {
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

	return ctx.Send("Привет, " + ctx.Sender().FirstName +
		"\n\nДоступные команды: \n" +
		"/addTask Добавление задачи с: " +
		"\n\t >названием" +
		"\n\t >описанием" +
		"\n\t >дедлайн задачи в формате дд-мм-гггг\n" +
		"/tasks получение списка своих задач \n" +
		"/deleteTask удаление задачи по номеру id")
}

func (bot *UpgradeBot) AddHandler(ctx telebot.Context) error {
	attempts := ctx.Args()

	existUser, err := bot.Users.FindOne(ctx.Chat().ID)

	if err != nil {
		log.Printf("Ошибка получения пользователя %v", err)
	}

	newTask := models.Task{
		Title:       attempts[0],
		Description: attempts[1],
		EndDate:     attempts[2],
		UserID:      existUser.ID,
	}

	err = bot.Tasks.Create(newTask)

	if err != nil {
		log.Printf("Ошибка создания задачи %v", err)
	}

	resStr := fmt.Sprintf(
		`Задача создана: 
			%s
			%s
			до %s`, newTask.Title, newTask.Description, newTask.EndDate,
	)
	return ctx.Send(resStr)
}

func (bot *UpgradeBot) AllTaskHandler(ctx telebot.Context) error {

	existUser, err := bot.Users.FindOne(ctx.Chat().ID)

	if err != nil {
		log.Printf("Ошибка получения пользователя %v", err)
	}

	tasks, err := bot.Tasks.GetAll(existUser.ID)

	if err != nil {
		log.Printf("Ошибка получения задачи %v", err)
	}

	list := make([]string, 0)

	for _, data := range tasks {
		str := fmt.Sprintf(
			`ID:  %d
					Задача:  %s
					Описание:  %s
					Дедлайн:  %s`,
			data.ID,
			data.Title,
			data.Description,
			data.EndDate,
		)
		list = append(list, str)
	}

	resStr := strings.Join(list, "\n\n")

	return ctx.Send(resStr)
}

func (bot *UpgradeBot) DeleteHandler(ctx telebot.Context) error {
	attempts := ctx.Args()

	if len(attempts) == 0 {
		return ctx.Send("Введите задачи")
	}

	if len(attempts) > 1 {
		return ctx.Send("Задача не найдено " +
			"введите по одной")
	}

	taskID, err := strconv.Atoi(attempts[0])

	if err != nil {
		log.Printf("Ошибка при ввода ID %v", err)
	}

	err = bot.Tasks.DeleteById(taskID)

	if err != nil {
		log.Printf("Ошибка при удалении задачи %v", err)
	}

	resStr := fmt.Sprintf("Задача №%d удалена", taskID)

	return ctx.Send(resStr)
}

func (bot *UpgradeBot) GameHandler(ctx telebot.Context) error {
	return ctx.Send("Сыграем в камень-ножницы-бумага " +
		"Введи твой вариант в формате /try камень")
}

func (bot *UpgradeBot) TryHandler(ctx telebot.Context) error {
	attempts := ctx.Args()

	if len(attempts) == 0 {
		return ctx.Send("Вы не ввели ваш вариант")
	}

	if len(attempts) > 1 {
		return ctx.Send("Вы ввели больше одного варианта")
	}

	try := strings.ToLower(attempts[0])
	botTry := gameItems[rand.Intn(3)]

	if botTry == "камень" {
		switch try {
		case "ножницы":
			ctx.Send(loseSticker)
			ctx.Send("🪨")
			return ctx.Send("Камень! Ты проиграл!")
		case "бумага":
			ctx.Send(winSticker)
			ctx.Send("🪨")
			return ctx.Send("Камень! Ты выиграл!")
		}
	}

	if botTry == "ножницы" {
		switch try {
		case "бумага":
			ctx.Send(loseSticker)
			ctx.Send("✂")
			return ctx.Send("Ножницы! Ты проиграл!")
		case "камень":
			ctx.Send(winSticker)
			ctx.Send("✂")
			return ctx.Send("Ножницы! Ты выиграл!")
		}
	}

	if botTry == "бумага" {
		switch try {
		case "ножницы":
			ctx.Send(winSticker)
			ctx.Send("📃")
			return ctx.Send("Бумага! Ты выиграл!")
		case "камень":
			ctx.Send(loseSticker)
			ctx.Send("📃")
			return ctx.Send("Бумага! Ты проиграл!")
		}
	}

	if botTry == try {
		return ctx.Send("Ничья!")
	}

	return ctx.Send("Кажется вы ввели неверный вариант!")
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
