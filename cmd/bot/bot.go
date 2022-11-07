package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"upgrade/internal/models"
	"upgrade/internal/repository"

	"gopkg.in/telebot.v3"
)

type TaskBot struct {
	Bot  *telebot.Bot
	Repo *repository.Repository
}

func InitBot(token string) *telebot.Bot {
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatalf("Ошибка при инициализации бота %v", err)
	}

	return bot
}

func (bot *TaskBot) InitCommands() {
	bot.Bot.Handle("/start", bot.start)
	bot.Bot.Handle("/addTask", bot.addTask)
	bot.Bot.Handle("/tasks", bot.tasks)
}

func (bot *TaskBot) start(ctx telebot.Context) error {
	newUser := models.User{
		Name:       ctx.Sender().Username,
		TelegramId: ctx.Sender().ID,
		FirstName:  ctx.Sender().FirstName,
		LastName:   ctx.Sender().LastName,
		ChatId:     ctx.Chat().ID,
	}

	user, err := bot.Repo.FindUser(ctx.Sender().ID)
	if err != nil {
		log.Printf("Ошибка получения пользователя %s", err)
	}

	if user == nil {
		if err := bot.Repo.CreateUser(newUser); err != nil {
			log.Printf("Ошибка создания пользователя %s", err)
		}
	}

	msg := fmt.Sprintf("Привет, %s", ctx.Sender().FirstName)

	return ctx.Send(msg)
}

func (bot *TaskBot) addTask(ctx telebot.Context) error {
	const (
		countTaskArgs = 3
		usageMsg      = "usage: /addTask название; описание; дедлайн;"
	)

	args := ctx.Args()
	if len(args) < countTaskArgs {
		msg := fmt.Sprintf("Должно быть %d аргумента для создания задачи\n", countTaskArgs)

		return ctx.Send(msg + usageMsg)
	}

	args = strings.Split(strings.Join(args, " "), ";")

	user, err := bot.Repo.FindUser(ctx.Sender().ID)
	if err != nil {
		log.Printf("Ошибка при добавлении задачи: %v", err)
		return ctx.Send("Не удалось добавить задачу")
	}

	task := models.Task{
		Title:       args[0],
		Description: args[1],
		EndDate:     args[2],
		UserId:      user.ID,
	}

	if err := bot.Repo.CreateTask(task); err != nil {
		log.Printf("Ошибка при добавлении задачи: %v", err)
		return ctx.Send("Не удалось добавить задачу")
	}

	return ctx.Send("Запись создана")
}

func (bot *TaskBot) tasks(ctx telebot.Context) error {
	user, err := bot.Repo.FindUser(ctx.Sender().ID)
	if err != nil {
		log.Printf("Ошибка при нахождении всех задач: %s", err)
		return ctx.Send("Не получилось найти задачи")
	}

	tasks, err := bot.Repo.AllTasks(user)
	if err != nil {
		log.Printf("Ошибка при нахождении всех задач: %s", err)
		return ctx.Send("Не получилось найти задачи")
	}

	if len(tasks) == 0 {
		return ctx.Send("У вас нет задач")
	}

	var msg string
	for i, task := range tasks {
		msg += fmt.Sprintf(
			"%d) Название: %s\nОписание: %s\nНужно сделать до: %s\n\n",
			i+1,
			task.Title,
			task.Description,
			task.EndDate,
		)
	}

	return ctx.Send(msg)
}
