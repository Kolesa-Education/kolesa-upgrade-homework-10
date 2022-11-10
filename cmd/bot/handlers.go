package bot

import (
	"gopkg.in/telebot.v3"
	"log"
	"strconv"
	"strings"
	"telegrambot/internal/models"
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

func (bot *UpgradeBot) TaskHandler(ctx telebot.Context) error {

	return ctx.Send("Введите задачу строго в таком порядке: /task Название ; Описание ; Дата окончания")

}

func (bot *UpgradeBot) AddTaskHandler(ctx telebot.Context) error {

	tempTask := ctx.Text()

	existUser, err := bot.Users.FindOne(ctx.Chat().ID)
	if err != nil {
		log.Printf("Ошибка получения пользователя %v", err)
	}
	task := strings.Split(tempTask, ";")
	newTask := models.Task{
		Title:       task[0][5:],
		Description: task[1],
		EndDate:     task[2],
		UserId:      existUser.ID,
	}
	err = bot.Tasks.Create(newTask)
	if err != nil {
		log.Printf("Ошибка создания задания %v", err)
	}

	return nil

}

func (bot *UpgradeBot) AllTaskHandler(ctx telebot.Context) error {
	user, err := bot.Users.FindOne(ctx.Chat().ID)
	if err != nil {
		log.Printf("Пользователь не зарегистрирован %v", err)
	}
	var tasks []models.Task
	err = bot.Users.Db.Model(&user).Association("Tasks").Find(&tasks)
	if err != nil {
		log.Printf("Ошибка при поиске задач %v", err)
	}
	if len(tasks) == 0 {
		return ctx.Send("У вас нет задач")
	}
	msg := ""
	for _, task := range tasks {
		msg += "Title:" + task.Title + "\n" + "Description:" + task.Description + "\n" + "Deadline:" + task.EndDate + "\n\n"
	}
	return ctx.Send(msg)
}

func (bot *UpgradeBot) DeleteTaskHandler(ctx telebot.Context) error {
	tempNum := ctx.Data()
	idDelete, err := strconv.Atoi(tempNum)
	if err != nil {
		log.Printf("Ошибка при считывании id")
		return ctx.Send("Пожалуйста повторите действие")
	}

	err = bot.Tasks.Delete(idDelete, "user_id")
	if err != nil {
		log.Printf("Ошибка при удалении записи %err", err)
		return ctx.Send("Произошла ошибка при удалении записи")
	}

	return nil
}
