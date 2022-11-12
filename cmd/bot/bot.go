package bot

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"kolesa-upgrade-hw-10/internal/models"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type UpgradeBot struct {
	Bot   *telebot.Bot
	Users *models.UserModel
	Tasks *models.TaskModel
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

func (bot *UpgradeBot) CreateTaskHandler(ctx telebot.Context) error {
	data := ctx.Data()
	task := strings.Split(data, ";")

	if len(task) != 3 {
		return ctx.Send("Некорректные данные для создания задачи, подробнее смотрите - /help")
	}

	endDate, er := time.Parse("02.01.2006", strings.TrimSpace(task[2]))

	if er != nil {
		return ctx.Send("Неверно указан формат даты, подробнее смотрите - /help")
	}

	user, err := bot.Users.FindOne(ctx.Sender().ID)
	newTask := models.Task{}

	if err == nil {
		newTask = models.Task{
			Title:       task[0],
			Description: task[1],
			EndDate:     endDate,
			UserId:      user.ID,
		}
	}

	err = bot.Tasks.Create(newTask)
	if err != nil {
		log.Printf("Ошибка создания задания %v", err)
		return ctx.Send("Ошибка создания задания")
	}

	return ctx.Send("Задание успешно добавлено")
}

func (bot *UpgradeBot) ShowTasksHandler(ctx telebot.Context) error {
	data, _ := os.ReadFile("messLayouts/tasks.txt")
	mess := ""
	tasks, _ := bot.Users.GetAllTasks()

	if len(tasks) == 0 {
		return ctx.Send("Список заданий пуст")
	}

	for _, task := range tasks {
		mess += fmt.Sprintf(
			string(data),
			task.Title,
			task.ID,
			task.Description,
			task.EndDate.Format("02.01.2006"),
		)
	}

	return ctx.Send(mess)
}

func (bot *UpgradeBot) DeleteTaskHandler(ctx telebot.Context) error {
	req := ctx.Args()
	mess := "нужно указать ID задания, вот так /deletetask 1"

	if len(req) == 0 {
		return ctx.Send(mess)
	}

	taskId := req[0]

	if id, err := strconv.Atoi(taskId); err == nil {
		mess = "Задача успешно удалена"
		err := bot.Tasks.DeleteTask(int64(id))

		if err != nil {
			mess = "Задача с указанным ID не существует"
		}
	}

	return ctx.Send(mess)
}

func (bot *UpgradeBot) HelpHandler(ctx telebot.Context) error {
	data, err := os.ReadFile("messLayouts/help.txt")
	mess := "Что-то пошло не так, попробуйте еще раз"

	if err == nil {
		mess = string(data)
	}

	return ctx.Send(mess)
}
