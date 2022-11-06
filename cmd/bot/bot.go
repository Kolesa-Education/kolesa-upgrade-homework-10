package bot

import (
	"log"
	"time"
	"upgrade/internal/models"

	"gopkg.in/telebot.v3"
	"fmt"
	"strconv"
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

func (bot *UpgradeBot) TaskHandler(ctx telebot.Context) error {
	return ctx.Send("Напиши название, описание, дату окончания задачи " +
		"/task английский, дз, 2022-12-08")
}

func (bot *UpgradeBot) WriteTaskHandler(ctx telebot.Context) error {
	attempts := ctx.Args()
	newTask := models.Task{
		Title:	attempts[0],
		Description:	attempts[1],
		End_date:	attempts[2],
		TelegramId :	ctx.Chat().ID,
	}

	err := bot.Tasks.Create(newTask)

	if err != nil {
		log.Printf("Ошибка создания задачи %v", err)
	}

	return ctx.Send("Задача создана, " + ctx.Sender().FirstName)
}

func (bot *UpgradeBot) ListTaskHandler(ctx telebot.Context) error {
	users, err := bot.Users.GetAll()
	if err != nil {
		return ctx.Send("Ошибка", err.Error())
	}
	var listtask string

	for _, user := range users {
		if user.ChatId == ctx.Chat().ID {
			for _, task := range user.Tasks {
				listtask += fmt.Sprintf("Название: %s\nОписание: %s\nДедлайн: %s", task.Title, task.Description, task.End_date)
			}
		}
	}

	if listtask == "" {
		return ctx.Send("У вас нет задач")
	}
	return ctx.Send(listtask)
}

func (bot *UpgradeBot) DeleteTaskHandler(ctx telebot.Context) error {
	arg := ctx.Args()

	if len(arg) == 0 {
		return ctx.Send("Вы не ввели id")
	}

	if len(arg) > 1 {
		return ctx.Send("Вы ввели больше одного варианта")
	}

	id, err := strconv.Atoi(arg[0])

	if err != nil {
		return ctx.Send("Вы ввели некорректный id Задание")
	}

	err = bot.Tasks.DeleteTask(id)
	if err != nil {
		return ctx.Send("Задание не удалено")
	}
	return ctx.Send("Задание удалено")
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