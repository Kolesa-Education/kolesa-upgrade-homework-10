package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"taskbot/internal/models"
	"time"

	"gopkg.in/telebot.v3"
)

type TaskBot struct {
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

func (bot *TaskBot) StartHandler(ctx telebot.Context) error {
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

func (bot *TaskBot) TaskRuleHandler(ctx telebot.Context) error {
	return ctx.Send("Сохраните задачи в бот " +
		"в формате: \n\"/addtask {Название};{Описание};{дд.мм.гггг(Дата окончания)}\"\n\n" +
		"Для удаления задачи по id пишите команду в формате: \n\"/deletetask {id задачи}\"\n\n" +
		"/tasks - Выдача всех задач пользователя\n\n")
}

func (bot *TaskBot) TaskHandler(ctx telebot.Context) error {
	taskmsg := ctx.Message().Text

	log.Print(taskmsg)

	if len(taskmsg) < 9 {
		return ctx.Send("Вы не ввели вашу задачу (помощь - /help)")
	}

	args := strings.Split(taskmsg[8:], ";")

	if len(args) > 3 {
		return ctx.Send("Вы ввели больше трех параметров (помощь - /help)")
	}
	if len(args) < 3 {
		return ctx.Send("Вы ввели недостаточно параметров (помощь - /help)")
	}

	title := args[0]
	description := args[1]
	end_date := args[2]
	fmt.Println(end_date)
	date, error := time.Parse("02.01.2006", strings.TrimSpace(end_date))

	if error != nil {
		fmt.Println(error)
	}

	newTask := models.Task{
		Title:       title,
		Description: description,
		EndDate:     date,
		UserId:      ctx.Chat().ID,
	}

	err := bot.Tasks.Create(newTask)

	if err != nil {
		log.Printf("Ошибка создания задачи %v", err)
	}

	return ctx.Send("Задача создана!")

}

func (bot *TaskBot) AllTasksHandler(ctx telebot.Context) error {
	user_id := ctx.Chat().ID

	tasks, err := bot.Tasks.FindTasks(user_id)
	if err != nil {
		log.Printf("Ошибка поиска задач %v", err)
	}

	result := ""

	for _, task := range tasks {
		result += "Id: " + strconv.FormatUint(uint64(task.ID), 10) + "\n"
		result += "Task: " + task.Title + "\n"
		result += "Description: " + task.Description + "\n"
		result += "End date: " + task.EndDate.Format("02.01.2006") + "\n\n"
	}

	return ctx.Send(result)
}

func (bot *TaskBot) DeleteTaskHandler(ctx telebot.Context) error {
	taskToDelete := ctx.Args()

	if len(taskToDelete) == 0 {
		return ctx.Send("Вы не ввели id задачи для удаления (помощь - /help)")
	}

	if len(taskToDelete) > 1 {
		return ctx.Send("Вы ввели больше одной задачи для удаления (помощь - /help)")
	}

	taskId := strings.ToLower(taskToDelete[0])
	taskIdInt, error := strconv.ParseInt(taskId, 0, 64)

	if error != nil {
		fmt.Println(error)
	}

	err := bot.Tasks.DeleteTask(taskIdInt)

	if err != nil {
		log.Printf("Ошибка удаления задачи %v", err)
	}

	return ctx.Send("Задача удалена!")
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
