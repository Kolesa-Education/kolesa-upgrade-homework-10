package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"telebot/internal/models"
	"time"

	"gopkg.in/telebot.v3"
)

type TeleBot struct {
	Bot   *telebot.Bot
	Users *models.UserModel
	Tasks *models.TaskModel
}

func (bot *TeleBot) StartHandler(ctx telebot.Context) error {
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

func (bot *TeleBot) HelpHandler(ctx telebot.Context) error {
	return ctx.Send("Можешь добавить новую задачу в формате:\n" +
		"/addTask Название задачи; Описание задачи; Дату сдачи задачи \n" +
		"/tasks просмотреть названия задач \n" +
		"/deleteTask {taskId} удалить задачу")
}

func splitTask(args *[]string) []string {
	argString := strings.Join(*args, " ")

	return strings.Split(argString, ";")
}

func (bot *TeleBot) AddTaskHandler(ctx telebot.Context) error {

	task := ctx.Args()
	task = splitTask(&task)

	newTask := models.Task{
		Title:       task[0],
		Description: task[1],
		EndDate:     task[2],
		UserID:      ctx.Sender().ID,
	}

	if len(task) == 0 {
		log.Println("Вы не добавили задачу")
	}

	if len(task) != 3 {
		log.Println("Неправильный формат задачи")
		ctx.Send("Неправильный формат задачи")
	}

	err := bot.Tasks.Create(newTask)

	if err != nil {
		log.Printf("Ошибка при создания задачи %v", err)
	}

	return ctx.Send("Задача " + newTask.Title + " создана")
}

func (bot *TeleBot) TasksHandler(ctx telebot.Context) error {

	userId := ctx.Sender().ID

	tasks, err := bot.Tasks.GetAll(userId)

	if err != nil {
		log.Printf("Ошибка при чтении задач %v", err)
		return ctx.Send("Ошибка при чтении задач")
	}

	result := ""

	for _, task := range tasks {

		result += fmt.Sprintf("Title: %s\nDescription: %s\nEndDate: %s\n\n", task.Title, task.Description, task.EndDate)
	}

	return ctx.Send("Список задач: \n" + result)
}

func (bot *TeleBot) DeleteTaskHandler(ctx telebot.Context) error {
	task := ctx.Args()

	taskId, err := strconv.Atoi(task[0])

	if err != nil {
		log.Printf("Ошибка при удалении задачи %v", err)
	}

	err = bot.Tasks.DeleteTask(taskId)

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
