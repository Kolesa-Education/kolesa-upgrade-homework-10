package bot

import (
	"log"
	"time"
	"strings"
	"strconv"
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

func (bot *UpgradeBot) AllTasksHandler(ctx telebot.Context) error {
	telegram_id := ctx.Chat().ID

	tasks, err := bot.Tasks.GetAll(telegram_id);

	if err != nil {
		log.Printf("Ошибка при сборе задании %v", err)
		return ctx.Send("Ошибка при сборе задании")
	}

	result := ""

	for _, task := range tasks {
		result += "ID: " + strconv.FormatUint(uint64(task.ID), 10) + "\n"
		result += "TASK: " + task.Title + "\n"
		result += "DESCRIPTION: " + task.Description + "\n"
		result += "END DATE: " + task.EndDate + "\n"
	}

	return ctx.Send(result)
}

func (bot *UpgradeBot) AddTaskHandler(ctx telebot.Context) error {
	txt := ctx.Message().Text

	log.Print(txt)

	txt = txt[9:]

	args := strings.Split(txt, ",")

	title := args[0]
	description := args[1]
	enddate := args[2]

	newTask := models.Task{
		Title:			title,
		Description: 	description,
		EndDate:		enddate,
		TelegramID:		ctx.Chat().ID,
	}

	err := bot.Tasks.Create(newTask);

	if err != nil {
		log.Printf("Ошибка создания задания %v", err)
		return ctx.Send("Ошибка создания задания")
	}

	return ctx.Send("Успешно создано!")
}

func (bot *UpgradeBot) DeleteTaskHandler(ctx telebot.Context) error {
	attempts := ctx.Args()

	if len(attempts) == 0 {
		return ctx.Send("Вы не ввели ваш вариант")
	}

	if len(attempts) > 1 {
		return ctx.Send("Вы ввели больше одного варианта")
	}

	id, err := strconv.Atoi(attempts[0])

	if err != nil {
		log.Printf("Ошибка конверсии аргумента в число %v", err)
		return ctx.Send("Ошибка конверсии аргумента в число")
	}

	er := bot.Tasks.Delete(id)

	if er != nil {
		log.Printf("Ошибка удаления задания %v", er)
		return ctx.Send("Ошибка удаления задания")
	}

	return ctx.Send("Успешно удалено!")
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