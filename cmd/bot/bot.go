package bot

import (
	"github.com/Kolesa-Education/kolesa-upgrade-homework-10/internal/models"
	"github.com/jedib0t/go-pretty/v6/table"
	"gopkg.in/telebot.v3"
	"log"
	"strconv"
	"strings"
	"time"
)

type Bot struct {
	Bot   *telebot.Bot
	Users *models.UserModel
	Tasks *models.TasksModel
}

func (bot *Bot) StartHandler(ctx telebot.Context) error {
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

	return ctx.Send("Привет" + ctx.Sender().FirstName)
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

func (bot *Bot) AddTaskHandler(ctx telebot.Context) error {
	attempts := ctx.Args()
	if len(attempts) != 3 {
		return ctx.Send("Не правильное количество аргументов, введите в формате /addTask Зоголовок; Описание; Дата;")
	}
	user, err := bot.Users.FindOne(ctx.Chat().ID)
	if err != nil {
		return ctx.Send("Server Error")
	}
	attempts = strings.Split(strings.Join(attempts, " "), ";")
	task := models.Tasks{
		Title:       attempts[0],
		Description: attempts[1],
		EndDate:     attempts[2],
		UserId:      user.ID,
	}
	bot.Tasks.Create(task)
	return ctx.Send("Вы добавили задание")
}

func (bot *Bot) TasksHandler(ctx telebot.Context) error {
	existUser, err := bot.Users.FindOne(ctx.Chat().ID)
	tasks, err := bot.Tasks.GetTasks(existUser.ID)
	if err != nil {
		return ctx.Send("Server Error")
	}
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Date", "Title", "Description"})
	for _, task := range tasks {
		t.AppendRows([]table.Row{
			{task.EndDate, task.Title, task.Description},
		})
	}
	t.AppendSeparator()
	return ctx.Send(t.RenderMarkdown())
}

func (bot *Bot) DeleteTaskHandler(ctx telebot.Context) error {
	attempts := ctx.Args()
	if len(attempts) != 1 {
		return ctx.Send("Неправильное количество аргументов, введите в формате /delete Task id")
	}
	id, err := strconv.Atoi(attempts[0])
	if err != nil {
		return ctx.Send("id должен быть числом")
	}
	err = bot.Tasks.DeleteTask(uint(id))
	if err != nil {
		ctx.Send("Несуществующий id")
	}
	return ctx.Send("Deleted")
}
