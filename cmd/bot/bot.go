package bot

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"log"
	"strconv"
	"strings"
	"time"
	"upgrade/internal/models"
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
	log.Printf("Бот успешно запущен")
	return ctx.Send("Привет, " + ctx.Sender().FirstName + "\nЭто бот планировщик задач\n" +
		"Используй команду /help")

}

func (bot *UpgradeBot) HelpHandler(ctx telebot.Context) error {
	log.Printf("Исползована команда /help")
	return ctx.Send("Список доступных комманд:\n\n" +
		"/new - Добавить новую задачу\n\n" +
		"/show - Показать список задач\n\n" +
		"/delete - Удалить задачу")
}

func (bot *UpgradeBot) NewTaskHandler(ctx telebot.Context) error {
	log.Printf("Исползована команда /new")
	text := ctx.Text()
	if text == "/new" {
		return ctx.Send("Введи заголовок, описание и дату задачи, " +
			"разделяя их дефисом. \nНапример: \n" +
			"/new Купить корм коту-Royal Canin-01.12.2022")
	}
	newText := strings.Replace(text, "/new ", "", -1)
	val := strings.Split(newText, "-")
	if len(val) != 3 {
		return ctx.Send("Введи заголовок, описание и дату задачи, " +
			"разделяя их дефисом. \nНапример: \n" +
			"/new Купить корм коту-Royal Canin-01.12.2022")
	}
	date, err := time.Parse("02.01.2006", val[2])
	if err != nil {
		log.Printf("Неправильный формат даты %v", err)
		return ctx.Send("Неправильный формат даты. \nПример: 01.12.2022")
	}
	newTask := models.Task{
		Title:       val[0],
		Description: val[1],
		UserID:      ctx.Sender().ID,
		EndDate:     date,
	}
	if err := bot.Tasks.CreateTask(newTask); err != nil {
		log.Printf("Ошибка добавления задачи %v", err)
		return ctx.Send("Ошибка добавления задачи")
	}
	log.Printf("Новая задача добавлена")
	return ctx.Send("Новая задача добавлена")
}

func (bot *UpgradeBot) ShowTaskHandler(ctx telebot.Context) error {
	log.Printf("Исползована команда /show")
	tasks, err := bot.Tasks.GetAll(ctx.Sender().ID)
	if err != nil {
		log.Fatalf("Ошибка получения списка задач %v", err)
	}
	var Tasks []string
	var task string
	for _, val := range tasks {
		task = fmt.Sprintf("ID: %v \nЗаголовок: %v \nОписание: %v \nДата: %v",
			val.ID, val.Title, val.Description, val.EndDate.Format("02.01.2006"))
		Tasks = append(Tasks, task)
	}
	result := strings.Join(Tasks, "\n\n")
	if result == "" {
		err = ctx.Send("Список задач пуст")
		if err != nil {
			return err
		}
	}
	return ctx.Send(result)
}

func (bot *UpgradeBot) DeleteTaskHandler(ctx telebot.Context) error {
	log.Printf("Исползована команда /delete")
	if ctx.Text() == "/delete" {
		return ctx.Send("Укажи id задачи для её удаления. \nНапример: \n/delete 2\n" +
			"Список задач можно посмотреть командой /show")
	}
	args := ctx.Args()
	deleteId, err := strconv.ParseInt(args[0], 0, 64)
	if err != nil {
		log.Printf("Ошибка удаления: ID не является целым числом %v", err)
		return ctx.Send("Ошибка удаления: ID не является целым числом")
	}
	if len(args) > 1 {
		return ctx.Send("Укажи один ID")
	}
	if err := bot.Tasks.DeleteTask(deleteId, ctx.Sender().ID); err != nil {
		log.Fatalf("Ошибка удаления задачи %v", err)
	}
	log.Printf("Задача %v успешно удалена ", deleteId)
	return ctx.Send("Задача успешно удалена")
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
