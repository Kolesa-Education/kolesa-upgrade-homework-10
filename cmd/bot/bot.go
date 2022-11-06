package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
	"upgrade/internal/models"
)

type UpgradeBot struct {
	Bot   *telebot.Bot
	Users *models.UserModel
	Tasks *models.TaskModel
}

func (b *UpgradeBot) StartHandler(ctx telebot.Context) error {
	newUser := models.User{
		Name:       ctx.Sender().Username,
		TelegramId: ctx.Chat().ID,
		FirstName:  ctx.Sender().FirstName,
		LastName:   ctx.Sender().LastName,
		ChatId:     ctx.Chat().ID,
	}

	user, err := b.Users.FindOne(ctx.Chat().ID)
	if err != nil {
		log.Printf("Пользователь %s зарегистрирован", ctx.Sender().Username)
	}

	if user == nil {
		err = b.Users.Create(newUser)
		if err != nil {
			log.Printf("Ошибка создания пользователя %v", err)
		}
	}
	return ctx.Send("Привет, " + ctx.Sender().FirstName)
}

func (b *UpgradeBot) AddTaskHandler(ctx telebot.Context) error {
	args := ctx.Args()
	if len(args) != 3 {
		msg := fmt.Sprintf("/addTask Название; Описание; Дедлайн;")
		ctx.Send(msg)
		return ctx.Send("Неверное количество аргументов для создания задачи. Нужно 3")
	}
	args = strings.Split(strings.Join(args, " "), ";")
	user, err := b.Users.FindOne(ctx.Chat().ID)
	if err != nil {
		log.Println(err)
		return nil
	}
	t := models.Task{
		Title:       args[0],
		Description: args[1],
		EndDate:     args[2],
		UserID:      user.ID,
	}
	b.Users.Db.Model(&user).Association("Tasks").Append([]models.Task{
		t,
	})
	return ctx.Send("Запись создана")
}

func (b *UpgradeBot) DeleteTaskHandler(ctx telebot.Context) error {
	id := ctx.Data()
	deleteId, err := strconv.Atoi(id)
	if err != nil || deleteId < 1 {
		ctx.Send("/deleteTask {id}")
		return ctx.Send("Неверный id")
	}
	user, err := b.Users.FindOne(ctx.Chat().ID)
	if err != nil {
		log.Println(err)
		return nil
	}
	t := models.Task{
		ID: uint(deleteId),
	}
	b.Users.Db.Model(&user).Association("Tasks").Delete([]models.Task{t})
	return ctx.Send("Запись удалена")
}

func (b *UpgradeBot) TasksHandler(ctx telebot.Context) error {
	user, err := b.Users.FindOne(ctx.Chat().ID)
	if err != nil {
		log.Printf("Пользователь %s зарегистрирован", ctx.Sender().Username)
	}
	var tasks []models.Task
	b.Users.Db.Model(&user).Association("Tasks").Find(&tasks)
	if len(tasks) == 0 {
		return ctx.Send("У вас нет задач")
	}
	for _, task := range tasks {
		msg := fmt.Sprintf("Задача %v\nНазвание: %s\nОписание: %s\nДедлайн: %v",
			task.ID,
			task.Title,
			task.Description,
			task.EndDate,
		)
		ctx.Send(msg)
	}
	return ctx.Send("ABOBa")
}

func InitBot(token string) *telebot.Bot {
	pref := telebot.Settings{
		Token: token,
		Poller: &telebot.LongPoller{
			Timeout: 10 * time.Second,
		},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatalf("Ошибка при инициализации бота %v", err)
	}

	return b
}
