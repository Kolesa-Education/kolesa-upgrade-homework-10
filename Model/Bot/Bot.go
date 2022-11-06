package Bot

import (
	"gopkg.in/telebot.v3"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	"upgrade/Model/Database"
)

type Bot struct {
	Database *Database.Database
	Bot      *telebot.Bot
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

func (bot *Bot) StartHandler(ctx telebot.Context) error {
	db := bot.Database.NewUser(ctx.Sender().ID, ctx.Sender().FirstName, ctx.Sender().LastName, ctx.Chat().ID)
	if db != nil {
		log.Println("Error creating user:", db.Error())
	}
	return ctx.Send("Привет, " + ctx.Sender().FirstName)
}

func (bot *Bot) ShowTasks(ctx telebot.Context) error {
	var result = "Tasks:\n"
	tasks, err := bot.Database.GetUserTasks(ctx.Sender().ID)
	if err != nil {
		return ctx.Send("Oops, something went wrong: Error getting tasks")
	}
	if len(tasks) == 0 {
		return ctx.Send("No tasks yet")
	}
	for _, task := range tasks {
		result += "ID: " + strconv.Itoa(task.Id) + "\n" +
			task.Title + "\n" +
			task.Description + "\n" +
			task.EndDate.Format("02.01.2006 15:04") + "\n\n"
	}
	return ctx.Send(result)
}

func (bot *Bot) NewTask(ctx telebot.Context) error {
	var args []string
	regex := regexp.MustCompile(`/\?`)
	args = regex.Split(strings.Replace(ctx.Text(), "/newTask ", "", 1), -1)
	if len(args) != 3 {
		return ctx.Send("Arguments count mismatch. Expected 3, got " + strconv.Itoa(len(args)))
	}
	if _, err := time.Parse("02.01.2006 15:04", args[2]); err != nil {
		return ctx.Send("Date must be in format: DD.MM.YYYY HH:mm")
	}
	err := bot.Database.NewTask(ctx.Sender().ID, args)
	if err != nil {
		return ctx.Send("Error creating new task")
	}
	return ctx.Send("New task:\n" +
		"Title:\t" + args[0] + "\n" +
		"Description:\t" + args[1] + "\n" +
		"End Time:\t" + args[2])
}

func (bot *Bot) DeleteTask(ctx telebot.Context) error {
	if len(ctx.Args()) != 1 {
		return ctx.Send("Arguments count mismatch. Expected 1, got " + strconv.Itoa(len(ctx.Args())))
	}
	taskId, ParseErr := strconv.Atoi(ctx.Args()[0])
	if ParseErr != nil {
		return ctx.Send("Define task ID to delete")
	}
	rowsAffected, err := bot.Database.DeleteTask(taskId, ctx.Sender().ID)
	if err != nil || rowsAffected == 0 {
		return ctx.Send("Unable to delete task with ID: " + ctx.Args()[0])
	}
	return ctx.Send("Task " + ctx.Args()[0] + " deleted")
}
