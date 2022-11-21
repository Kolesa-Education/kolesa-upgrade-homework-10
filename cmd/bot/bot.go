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

func (bot *UpgradeBot) AddTaskHandler(ctx telebot.Context) error {
	chatId := ctx.Chat().ID
	message := ctx.Message().Text
	message = message[8:]
	data := strings.Split(message, "|")

	log.Println(data)
	if len(data) < 3 {
		return ctx.Send("Неправильный формат задачи")
	}

	for i := 0; i < len(data); i++ {
		data[i] = strings.Trim(data[i], "\t \f \v")
	}

	title := data[0]
	description := data[1]
	endDate, err := time.Parse("02.06.2006", data[2])

	fmt.Print(endDate)

	if err != nil {
		return ctx.Send("Вы ввели неправильную дату. \nПример даты: 07.07.2022")
	}

	newTask := models.Task{
		Title:       title,
		Description: description,
		EndDate:     endDate,
		ChatId:      chatId,
	}

	err = bot.Tasks.Create(newTask)

	if err != nil {
		log.Printf("Ошибка создания пользователя %v", err)
		return ctx.Send("Ошибка создания пользователя!")

	}

	return ctx.Send("Задача добавлена!")
}

func (bot *UpgradeBot) AllTasksHandler(ctx telebot.Context) error {
	chatId := ctx.Chat().ID

	existTasks, err := bot.Tasks.FindAll(chatId)

	if err != nil {
		return ctx.Send("Ошибка при получений задаач!")
	}

	var (
		Tasks     []string
		task_info string
	)

	for _, el := range existTasks {
		task_info = fmt.Sprintf("ID: %d\nTitle: %s\nDescription: %s\nDeadline: %s", el.ID, el.Title, el.Description, el.EndDate.Format("01.02.2006"))
		Tasks = append(Tasks, task_info)
	}

	resString := strings.Join(Tasks, "\n\n")

	return ctx.Send(resString)
}

func (bot *UpgradeBot) DeleteTaskHandler(ctx telebot.Context) error {
	args := ctx.Args()
	taskId, err := strconv.ParseInt(args[0], 6, 12)

	if err != nil {
		return ctx.Send("Ошибка id")
	}

	err = bot.Tasks.Delete(taskId)

	if err != nil {
		return ctx.Send("Ошибка удаления id")
	}

	return ctx.Send("Удалено")

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
