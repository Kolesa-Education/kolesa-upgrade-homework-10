package bot

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"upgrade/internal/models"

	"gopkg.in/telebot.v3"
)

type UpgradeBot struct {
	Bot   *telebot.Bot
	Tasks *models.TaskModel
	Users *models.UserModel
}

var gameItems = [3]string{
	"камень",
	"ножницы",
	"бумага",
}

var winSticker = &telebot.Sticker{
	File: telebot.File{
		FileID: "CAACAgIAAxkBAAEGMEZjVspD4JulorxoH7nIwco5PGoCsAACJwADr8ZRGpVmnh4Ye-0RKgQ",
	},
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
		TelegramId: ctx.Sender().ID,
		FirstName:  ctx.Sender().FirstName,
		LastName:   ctx.Sender().LastName,
		ChatId:     ctx.Chat().ID,
	}

	existUser, err := bot.Users.FindOne(ctx.Sender().ID)

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
	botTry := gameItems[rand.Intn(len(gameItems))]

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
		case "камень":
			ctx.Send(winSticker)
			ctx.Send("✂️")
			return ctx.Send("Ножницы! Ты выиграл!")
		case "бумага":
			ctx.Send(loseSticker)
			ctx.Send("✂️")
			return ctx.Send("Ножницы! Ты проиграл!")
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

func (bot *UpgradeBot) AddTaskHandler(ctx telebot.Context) error {
	return ctx.Send("Заполните задачу в формате : " +
		"/newtask Заголовок//Описание//Дедлайн(чч:мм дд.мм.гг)")
}

func (bot *UpgradeBot) NewTaskHandler(ctx telebot.Context) error {

	input := ctx.Args()

	if len(input) == 0 {
		return ctx.Send("Вы не заполнили задачу!")
	}

	args := strings.Join(input, " ")
	args = strings.Replace(args, "/newtask ", "", 1)
	argsArr := strings.Split(args, "//")

	if len(argsArr) > 3 {
		return ctx.Send("Вы ввели лешние разделители '//'")
	}

	title := argsArr[0]
	description := argsArr[1]
	endDate := argsArr[2]
	existUser, err := bot.Users.FindOne(ctx.Sender().ID)
	userId := existUser.ID

	newTask := models.Task{
		Title:       title,
		Description: description,
		EndDate:     endDate,
		UserId:      int64(userId),
	}

	existTask, err := bot.Tasks.FindSame(0, title, description)

	if err != nil {
		log.Printf("Ошибка получения задачи %v", err)
	}

	if existTask == nil {
		err := bot.Tasks.Create(newTask)

		if err != nil {
			log.Printf("Ошибка создания задачи %v", err)
		}
	}

	return ctx.Send("Новая задача: " + title + "\n " + description + "\n Дедлайн: " + endDate)

}

func (bot *UpgradeBot) TasksHandler(ctx telebot.Context) error {
	existUser, err := bot.Users.FindOne(ctx.Sender().ID)
	if err != nil {
		return ctx.Send("У вас ещё нет задач, сипользуйте команду /start звтем /addtask")
	}
	userId := existUser.ID
	userTasks, _ := bot.Tasks.GetAllByUserId(int64(userId))
	var tasks []string

	for _, task := range userTasks {
		taskId := strconv.Itoa(task.ID)
		tasks = append(tasks, "id: "+taskId+"\n"+task.Title+"\n"+task.Description+"\n"+task.EndDate+"\n")
	}
	result := strings.Join(tasks, "\n")
	return ctx.Send(result)
}

func (bot *UpgradeBot) DeleteTaskHandler(ctx telebot.Context) error {
	input := ctx.Args()

	if len(input) > 1 {
		return ctx.Send("Вы ввели слишком много аргументов!")
	}

	taskId, err := strconv.Atoi(input[0])

	if err != nil {
		return ctx.Send("id задачи неверное!")
	}

	existTask, err := bot.Tasks.FindSame(taskId, "", "")

	if err != nil {
		return ctx.Send("Неудалось найти задачу!")
	}

	err = bot.Tasks.DropTask(*existTask)
	if err != nil {
		return ctx.Send("Не удалось удалить задачу, попробуйте ещё раз!")
	}

	return ctx.Send("Задача удалена")
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
