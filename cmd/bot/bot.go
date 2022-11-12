package bot

import (
	"gopkg.in/telebot.v3"
	"kolesa-upgrade-hw-10/internal/models"
	"log"
	"os"
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
	task := ctx.Args()
	if len(task) != 3 {
		return ctx.Send("Некорректные данные для создания задачи, подробнее смотрите - /help")
	}
	endDate, er := time.Parse("02.01.2006", task[2])

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

func (bot *UpgradeBot) HelpHandler(ctx telebot.Context) error {
	data, err := os.ReadFile("messLayouts/help.txt")
	mess := "Какие-то неопладки, попробуйте еще раз"

	if err == nil {
		mess = string(data)
	}

	return ctx.Send(mess)
}
