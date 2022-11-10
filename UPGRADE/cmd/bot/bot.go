package bot

import (
	"fmt"
	"log"

	"strconv"
	"strings"
	"time"
	"upgrade/internal/models"

	"gopkg.in/telebot.v3"
)

type UpgradeBot struct {
	Bot   *telebot.Bot
	Users *models.UserModel
	Tasks *models.TaskModel
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
		log.Printf("создай юзера %v", err)
	}

	if existUser == nil {
		err := bot.Users.Create(newUser)

		if err != nil {
			log.Printf("Ошибка создания пользователя %v", err)
		}
	}

	return ctx.Send("Привет, " + ctx.Sender().FirstName)
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
func (bot *UpgradeBot) AddTask(ctx telebot.Context) error {
	task := ctx.Text()
	value := strings.Replace(task, "/addTask ", "", -1)
	values := strings.Split(value, "/")

	t, _ := time.Parse("2006-01-02 15:04:05", values[2])
	newTask := models.Task{
		Task:        values[0],
		Description: values[1],
		TelegramId:  ctx.Chat().ID,
		End_date:    t,
	}
	log.Println(newTask.End_date)
	err := bot.Tasks.Create(newTask)

	if err != nil {
		log.Printf("Ошибка создания Task %v", err)
		return ctx.Send("Task не создана")
	}
	return ctx.Send("Task создана")
}
func (bot *UpgradeBot) ShowTasks(ctx telebot.Context) error {
	var result = "Задачки:\n"
	tasks, err := bot.Tasks.AllTask(ctx.Chat().ID)

	if err != nil {
		log.Printf("Ошибка получения задач %v", err)
		return ctx.Send("Ошибка получения задач")
	}

	for _, task := range tasks {
		result += task.Task + "\n" + task.Description + "\n" + fmt.Sprintln(task.End_date.Format("2006-01-02 15:04:05")) + "\n"
	}
	return ctx.Send(result)
}

func (bot *UpgradeBot) TaskDel(ctx telebot.Context) error {
	arg := ctx.Args()
	taskId, err := strconv.Atoi(arg[0])

	if err != nil {
		//panic(err)
		log.Printf("Ошибка преобразования %v", err)
		return ctx.Send("Введите число")
	}

	if err := bot.Tasks.DeleteTask(taskId, ctx.Chat().ID); err != nil {
		log.Printf("Ошибка удаления %v", err)
		return ctx.Send("Ошибка удаления")
	}

	return ctx.Send("task удалена")
}
