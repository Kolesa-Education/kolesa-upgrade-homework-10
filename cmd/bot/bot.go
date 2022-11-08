package bot

import (
	"fmt"
	"log"
	"strings"
	"telegramTaskBot/internal/models"
	"time"

	"gopkg.in/telebot.v3"
)

type UpgradeBot struct {
	Bot   *telebot.Bot
	Users *models.UserModel
	Tasks *models.TaskModel
}

func (bot *UpgradeBot) AddHandler(ctx telebot.Context) error {
	task := ctx.Text()
	value := strings.Replace(task, "/addTask ", "", -1)

	vals := strings.Split(value, ",")

	date, err1 := time.Parse("02.01.2006 15:04", vals[2])
	if err1 != nil {
		fmt.Println(len(vals), vals, date)

	}
	newTask := models.Task{
		Title:       vals[0],
		Description: vals[1],
		UserId:      ctx.Sender().ID,
		EndDate:     date,
	}
	fmt.Println(date)
	if err := bot.Tasks.Create(newTask); err != nil {
		log.Printf("Ошибка создания задания %v", err)
		return ctx.Send("Ошибка создания задания")
	}
	return ctx.Send("Задание успешно добавлено")
}

func (bot *UpgradeBot) ShowHandler(ctx telebot.Context) error {

	tasks, err := bot.Tasks.FindAll(ctx.Sender().ID)

	if err != nil {
		log.Fatalf("Ошибка обработки задачи %v", err)
	}

	var strTasks []string

	for _, el := range tasks {
		str := fmt.Sprintf(`#: %d: %s: %s due: %s`, el.ID, el.Title, el.Description, el.EndDate)
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
	return ctx.Send("Привет, " + ctx.Sender().FirstName + "\n/addTask TaskName,Description,DueDate\n/tasks")

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
