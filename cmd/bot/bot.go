package bot

import (
	"bytes"
	"log"
	"strconv"
	"strings"
	"tasking/internal/models"
	"time"

	"gopkg.in/telebot.v3"
	"gorm.io/datatypes"
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

	return ctx.Send("Привет, " + strconv.Itoa(existUser.Id))
}

func (bot *UpgradeBot) AddTaskHandler(ctx telebot.Context) error {
	args := ctx.Text()

	if strings.Count(args, "\n") != 2 {
		return ctx.Send("Неверный формат ввода задания")
	}

	text := strings.Split(args, "\n")
	var title bytes.Buffer
	for i, value := range strings.Split(text[0], " ") {
		if i != 0 {
			title.WriteString(value)
			if i != len(text[0])-1 {
				title.WriteString(" ")
			}
		}
	}

	existUser, err := bot.Users.FindOne(ctx.Chat().ID)
	if err != nil {
		log.Printf("Ошибка получения пользователя %v", err)
		return ctx.Send("К сожалению, На данный момент мы не можем сохранить ваше задание")
	}
	user_id := existUser.Id

	format := "2006-01-02"
	dateParsing, _ := time.Parse(format, "2019-07-10")
	date := datatypes.Date(dateParsing)

	newTask := models.Task{
		Title:       title.String(),
		Description: text[1],
		EndDate:     date,
		UserId:      user_id,
	}

	tasking := bot.Tasks.Create(newTask)
	if tasking != nil {
		log.Printf("Ошибка создания пользователя %v", tasking)
	}
	return ctx.Send("Задание " + title.String() + "сохранено")
}

func (bot *UpgradeBot) TasksHandler(ctx telebot.Context) error {
	existUser, _ := bot.Users.FindOne(ctx.Chat().ID)

	tasks, err := bot.Tasks.GetTasks(existUser.Id)

	if err != nil {
		log.Println(err)
		return ctx.Send("К сожалению произошла ошибка")
	}
	result := ""

	for _, task := range tasks {
		y, m, d := time.Time(task.EndDate).Date()
		// log.Println(task.EndDate.Value())
		result = result + task.Title + "\n" + task.Description + "\n" + strconv.Itoa(y) + " " + m.String() + " " + strconv.Itoa(d) + "\n\n"
	}
	return ctx.Send(result)
}

func (bot *UpgradeBot) DeleteTaskHandler(ctx telebot.Context) error {
	args := ctx.Text()

	if len(args) < 8 {
		return ctx.Send("Неверный формат ввода удаления задания")
	}

	task_id, _ := strconv.Atoi(args[8:])
	err := bot.Tasks.Delete(task_id)
	if err != nil {
		log.Printf("Ошибка удаления задания %v", err)
		return ctx.Send("Произошла ошибка")
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
