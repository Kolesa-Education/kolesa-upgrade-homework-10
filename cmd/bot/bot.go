package bot

import (
	"fmt"
	"log"
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

//func (b *UpgradeBot) AddTaskHandler(ctx telebot.Context) error {
//	ctx.Reply("Введите название задачи")
//	fmt.Println(ctx.Data())
//	telewa
//	ctx.Reply("Введите описание задачи")
//	fmt.Println(ctx.Data())
//	return ctx.Send("done")
//}

func (b *UpgradeBot) TasksHandler(ctx telebot.Context) error {
	user, err := b.Users.FindOne(ctx.Chat().ID)
	if err != nil {
		log.Printf("Пользователь %s зарегистрирован", ctx.Sender().Username)
	}
	var tasks []models.Task
	b.Tasks.Db.Model(&user).Association("Tasks").Find(&tasks)
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
