package bot

import (
	"firstTelegramBot/internal/models"
	"fmt"
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
	return ctx.Send("Привет, " + ctx.Sender().FirstName +
		"\n\nДоступные команды: \n" +
		"/addTask Заголовок; Описание; Дедлайн задачи в формате дд-мм-гггг\n" +
		"- добавление задачи в список (по шаблону выше) \n" +
		"/tasks получение списка своих задач \n" +
		"/deleteTask id - удаление задачи номер id")
}

func (bot *UpgradeBot) AddHandler(ctx telebot.Context) error {
	taskArgs := ctx.Args()
	refString := strings.Join(taskArgs, " ")
	taskArgs = strings.Split(refString, ";")

	existUser, err := bot.Users.FindOne(ctx.Chat().ID)
	if err != nil {
		log.Printf("Ошибка получения пользователя %v", err)
	}

	newTask := models.Task{
		Title:       strings.TrimSpace(taskArgs[0]),
		Description: strings.TrimSpace(taskArgs[1]),
		EndDate:     strings.TrimSpace(taskArgs[2]),
		UserID:      existUser.ID,
	}

	err = bot.Tasks.Create(newTask)

	if err != nil {
		log.Printf("Ошибка создания задачи %v", err)
	}

	resStr := fmt.Sprintf(
		`Создана задача: 
			%s
			%s
			до %s`, newTask.Title, newTask.Description, newTask.EndDate,
	)
	return ctx.Send(resStr)
}

func (bot *UpgradeBot) ShowHandler(ctx telebot.Context) error {
	existUser, err := bot.Users.FindOne(ctx.Chat().ID)
	if err != nil {
		log.Printf("Ошибка получения пользователя %v", err)
	}

	tasks, err := bot.Tasks.FindAll(existUser.ID)
	if err != nil {
		log.Printf("Ошибка получения задач %v", err)
	}

	var (
		taskSlice []string
		str       string
	)

	for _, taskItem := range tasks {
		str = fmt.Sprintf(
			`ID:  %d
					Задача:  %s
					Описание:  %s
					Дедлайн:  %s`,
			taskItem.ID,
			taskItem.Title,
			taskItem.Description,
			taskItem.EndDate,
		)
		taskSlice = append(taskSlice, str)
	}

	resString := strings.Join(taskSlice, "\n\n")

	return ctx.Send(resString)
}

func (bot *UpgradeBot) DeleteHandler(ctx telebot.Context) error {
	taskArg := ctx.Args()

	if len(taskArg) == 0 {
		return ctx.Send("Ошибка получения iD задачи")
	}

	if len(taskArg) > 1 {
		return ctx.Send("Ошибка получения iD задачи" +
			"Чтобы посмотреть iD нужной задачи, введите /tasks")
	}

	taskID, err := strconv.Atoi(taskArg[0])
	if err != nil {
		log.Printf("Ошибка удаления задачи. Неправильный iD %v", err)
	}

	err = bot.Tasks.DeleteOne(taskID)
	if err != nil {
		log.Printf("Ошибка удаления задачи %v", err)
	}

	return ctx.Send(fmt.Sprintf("Задача #%d удалена", taskID))
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
