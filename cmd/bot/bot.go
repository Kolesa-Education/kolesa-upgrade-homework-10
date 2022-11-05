package bot

import (
	"ZakirAvrora/homework-10/internal/models"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
)

var ErrInArguments = errors.New("invalid argument")

type UpgradeBot struct {
	Bot   *telebot.Bot
	Users *models.UserModel
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

	return ctx.Send("Привет, " + ctx.Sender().FirstName)
}

func (bot *UpgradeBot) GameHandler(ctx telebot.Context) error {
	return ctx.Send("Сыграем в камень-ножницы-бумага " +
		"Введи твой вариант в формате /try камень")
}

func (bot *UpgradeBot) AddTaskHandler(ctx telebot.Context) error {
	str := ctx.Data()

	if len(str) == 0 {
		return ctx.Send("Вы не ввели ваш таск")
	}

	args, err := Parse(str)
	if err != nil {
		return ctx.Send("Вы ввели неправильно")
	}

	task := models.Task{
		Title:       args[0],
		Description: args[1],
		EndDate:     args[2],
		TelegramId:  ctx.Chat().ID,
	}

	err = bot.Users.CreateTask(task)
	if err != nil {
		return ctx.Send("Ошибка", err.Error())
	}

	return ctx.Send("Успешно добавлена")
}

func (bot *UpgradeBot) TasksHandler(ctx telebot.Context) error {

	users, err := bot.Users.GetAll()
	if err != nil {
		return ctx.Send("Ошибка", err.Error())
	}
	var tasksMsg string

	for _, user := range users {
		if user.ChatId == ctx.Chat().ID {
			for _, task := range user.Tasks {
				tasksMsg += fmt.Sprintf("Загаловок: %s\nОписания: %s\nДедлайн: %s\n", task.Title, task.Description, task.EndDate)
			}
		}
	}

	if tasksMsg == "" {
		return ctx.Send("No task")
	}
	return ctx.Send(tasksMsg)
}

func (bot *UpgradeBot) DeleteTaskHandler(ctx telebot.Context) error {
	arg := ctx.Args()

	if len(arg) == 0 {
		return ctx.Send("Вы не ввели ваш вариант")
	}

	if len(arg) > 1 {
		return ctx.Send("Вы ввели больше одного варианта")
	}

	id, err := strconv.Atoi(arg[0])

	if err != nil {
		return ctx.Send("Вы ввели неправильный id таска")
	}

	err = bot.Users.DeleteTask(id)
	if err != nil {
		return ctx.Send("Cannot be deleted")
	}
	return ctx.Send("Deleted")
}

func Parse(str string) ([]string, error) {
	args := strings.Split(strings.TrimSpace(str), ":")

	if len(args) != 3 {
		return nil, ErrInArguments
	}
	return args, nil
}
