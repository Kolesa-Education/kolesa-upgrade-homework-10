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
	return ctx.Send("Напишите пожалуйста корректно: /task $Задание $Обьяснение $Конец")
}

func (bot *UpgradeBot) TaskHandler(ctx telebot.Context) error {
	tmpTask := ctx.Text()
	existUser, err := bot.Users.FindOne(ctx.Chat().ID)
	if err != nil {
		log.Printf("Ошибка получения пользователя %v", err)
	}

	task := strings.Split(tmpTask, "$")
	newTask := models.Task{
		Title:       task[1],
		Description: task[2],
		EndDate:     task[3],
		UserId:      existUser.ID,
	}
	err = bot.Tasks.Create(newTask)
	if err != nil {
		log.Printf("Ошибка создания пользователя %v", err)
	}
	return ctx.Send("Добавлена задача: " + newTask.Title)
}
func (bot *UpgradeBot) AllTasksHandler(ctx telebot.Context) error {
	existUser, err := bot.Users.FindOne(ctx.Chat().ID)
	if err != nil {
		log.Printf("Ошибка получения пользователя %v", err)
	}
	var tmpTask []models.Task
	err = bot.Users.Db.Model(&existUser).Association("Tasks").Find(&tmpTask)
	if err != nil {
		log.Printf("Ошибка при поиске задач %v", err)
	}
	if len(tmpTask) == 0 {
		return ctx.Send("У вас нет задач")
	}
	msg := ""
	for _, task := range tmpTask {
		msg += "Title:" + task.Title + "\n" + "Description:" + task.Description + "\n" + "Deadline:" + task.EndDate + "\n" + "\n"
	}
	return ctx.Send(msg)
}

func (bot *UpgradeBot) DeleteTasksHandler(ctx telebot.Context) error {
	tmpMsg := ctx.Data()
	num, err := strconv.Atoi(tmpMsg)
	if err != nil {
		return err
	}
	fmt.Println(num)
	existUser, err := bot.Users.FindOne(ctx.Chat().ID)
	if err != nil {
		log.Printf("Ошибка получения пользователя %v", err)
		return ctx.Send("Такого пользователя нет")
	}
	t := models.Task{
		ID: uint(num),
	}
	bot.Users.Db.Model(&existUser).Association("Tasks").Delete(&t)

	return nil
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
