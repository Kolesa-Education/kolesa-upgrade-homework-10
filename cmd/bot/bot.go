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

func (bot *UpgradeBot) DeleteHandler(ctx telebot.Context) error {

	attempts := ctx.Args()
	deleteId, _ := strconv.ParseInt(attempts[0], 0, 64)

	if len(attempts) == 0 {
		return ctx.Send("Вы не ввели ваш ИД")
	}

	if len(attempts) > 1 {
		return ctx.Send("Вы ввели больше одного ИД")
	}

	if err := bot.Tasks.DeleteTask(deleteId, ctx.Sender().ID); err != nil {
		log.Fatalf("Ошибка выполнения запроса пользователя %v", err)
	}

	return ctx.Send("Ваше задание удалено!")

}

func (bot *UpgradeBot) AddHandler(ctx telebot.Context) error {
	allText := ctx.Text()
	clearText := strings.Replace(allText, "/addTask ", "", -1)

	vals := strings.Split(clearText, ";")

	date, _ := time.Parse("02.01.2006 15:04", vals[2])
	newTask := models.Task{
		Title:       vals[0],
		Description: vals[1],
		UserId:      ctx.Sender().ID,
		EndDate:     date,
	}

	if err := bot.Tasks.Create(newTask); err != nil {
		log.Fatalf("Ошибка создания задания %v", err)
	}
	return ctx.Send("Новое задание успешно добавлено")
}

func (bot *UpgradeBot) ShowHandler(ctx telebot.Context) error {

	tasks, err := bot.Tasks.FindAll(ctx.Sender().ID)

	if err != nil {
		log.Fatalf("Ошибка обработки задачи %v", err)
	}

	var (
		strTasks []string
		str      string
	)

	for _, el := range tasks {
		str = fmt.Sprintf(`ИД: %d Задача: %s Описание: %s Срок Выполнения: %s`, el.ID, el.Title, el.Description, el.EndDate)
		strTasks = append(strTasks, str)
	}

	resString := strings.Join(strTasks, "\n")

	return ctx.Send(resString)
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

	return ctx.Send("Привет, " + ctx.Sender().FirstName + " /start - для начала работы /addTask - для добавления задания  /deleteTask - для удаления задания по ИД /tasks - для показа всех задания пользователя")
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
