package bot

import (
	"database/sql"
	"errors"
	"gobot/internal/models"
	"log"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
)

type UpgradeBot struct {
	Bot         *telebot.Bot
	Users       *models.UserModel
	Tasks       *models.TaskModel
	taskHandler map[int64]models.Task
}

const (
	datelayout = "2006-01-02"
)

func (bot *UpgradeBot) StartHandler(ctx telebot.Context) error {
	newUser := models.User{
		Name:       ctx.Sender().Username,
		TelegramId: ctx.Sender().ID,
		FirstName:  ctx.Sender().FirstName,
		LastName:   ctx.Sender().LastName,
		ChatId:     ctx.Chat().ID,
	}

	existUser, err := bot.Users.FindOne(ctx.Chat().ID)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("Ошибка получения пользователя %v", err)
	}

	if existUser == nil {
		err := bot.Users.Create(newUser)

		if err != nil {
			log.Printf("Ошибка создания пользователя %v", err)
		}
	}
	bot.taskHandler = make(map[int64]models.Task)

	return ctx.Send("Привет, " + ctx.Sender().FirstName)
}

func InitBot(token string) *telebot.Bot {
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)

	if err != nil {
		log.Fatalf("Bot initialization error %v", err)
	}

	return b
}

func (bot *UpgradeBot) AddTask(ctx telebot.Context) error {
	t, ok := bot.taskHandler[ctx.Sender().ID]
	if !ok {
		return ctx.Send("create a task before adding task itself.\nUse /addTaskTitle for title.\nUse /addTaskDescription for description.\nUse /addTaskEndDate for end date.")
	}
	if t.Title == "" {
		return ctx.Send("Title is not provided")
	}
	if t.Description == "" {
		return ctx.Send("Description is not provided")
	}

	if t.EndDate.IsZero() {
		return ctx.Send("End date is not provided")
	}

	t.TelegramID = ctx.Sender().ID
	if err := bot.Tasks.AddTask(t); err != nil {
		log.Printf("Ошибка создания задания %v", err)
	}

	delete(bot.taskHandler, ctx.Sender().ID)

	return nil
}

func (bot *UpgradeBot) AddTaskTitle(ctx telebot.Context) error {
	title := strings.Join(ctx.Args(), " ")
	t, ok := bot.taskHandler[ctx.Sender().ID]
	if ok {

		t.Title = title
		bot.taskHandler[ctx.Sender().ID] = t
		return nil
	}
	bot.taskHandler[ctx.Sender().ID] = models.Task{Title: title}
	return nil
}

func (bot *UpgradeBot) AddTaskDescription(ctx telebot.Context) error {
	desc := strings.Join(ctx.Args(), " ")
	t, ok := bot.taskHandler[ctx.Sender().ID]
	if ok {
		t.Description = desc
		bot.taskHandler[ctx.Sender().ID] = t
		return nil
	}
	bot.taskHandler[ctx.Sender().ID] = models.Task{Description: desc}
	return nil
}

func (bot *UpgradeBot) AddTaskEndDate(ctx telebot.Context) error {
	date, err := time.Parse(datelayout, ctx.Args()[0])
	if err != nil {
		return err
	}
	t, ok := bot.taskHandler[ctx.Sender().ID]
	if ok {
		t.EndDate = date
		bot.taskHandler[ctx.Sender().ID] = t
		return nil
	}
	bot.taskHandler[ctx.Sender().ID] = models.Task{EndDate: date}
	return nil
}

// func (bot *UpgradeBot) GetAll(ctx telebot.Context) error {

// }

// func (bot *UpgradeBot) DeleteTask(ctx telebot.Context) error {

// }
