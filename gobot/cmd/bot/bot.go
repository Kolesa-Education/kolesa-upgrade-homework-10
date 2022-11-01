package bot

import (
	"fmt"
	"gobot/internal/models"
	"log"
	"strconv"
	"strings"
	"time"

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
		TelegramId: ctx.Sender().ID,
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

	attempts := ctx.Text()
	truevalue := strings.Replace(attempts, "/addTask ", "", -1)
	fmt.Println(truevalue)
	values := strings.Split(truevalue, "/")
	for i, s := range values {
		fmt.Println(i, s)
	}
	intdate, _ := strconv.ParseInt(values[2], 10, 64)
	newTask := models.Task{
		Title:       values[0],
		Description: values[1],
		UserId:      ctx.Sender().ID,
		EndDate:     intdate,
	}
	log.Println(newTask.EndDate)
	err := bot.Tasks.Create(newTask)

	if err != nil {
		log.Printf("Ошибка создания пользователя %v", err)
	}
	return ctx.Send("GUT")
}
func (bot *UpgradeBot) AllTaskHandler(ctx telebot.Context) error {
	var task models.Task
	result := bot.Tasks.GetAll(ctx.Sender().ID)
	for result.Next() {
		bot.Tasks.Db.ScanRows(result, &task)
		log.Println(task.EndDate)
		ctx.Send(strconv.Itoa(int(task.ID)) + task.Title + task.Description + strconv.FormatInt(task.EndDate, 10))
	}
	return ctx.Send(result)
}

func (bot *UpgradeBot) DeleteTaskHandler(ctx telebot.Context) error {
	attempts := ctx.Args()

	if len(attempts) == 0 {
		return ctx.Send("Вы не ввели ваш вариант")
	}

	if len(attempts) > 1 {
		return ctx.Send("Вы ввели больше одного варианта")
	} else if len(attempts) == 1 {
		bot.Tasks.DeleteTask(attempts[0], ctx.Sender().ID)
	}

	return ctx.Send("DONE")
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
