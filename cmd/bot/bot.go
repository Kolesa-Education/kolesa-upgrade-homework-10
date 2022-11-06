package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"upgrade/cmd/task"
	"upgrade/internal/models"

	"gopkg.in/telebot.v3"
)

type TodoBot struct {
	Bot   *telebot.Bot
	Users *models.UserModel
	Tasks *models.TaskModel
}

func (bot *TodoBot) StartHandler(ctx telebot.Context) error {
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

func (bot *TodoBot) HelpHandler(ctx telebot.Context) error {
	return ctx.Send(
		"**Список комманд**\n\n" +
			"Чтобы добавить задачу, введи её в формате\n" +
			"/add Заголовок; Описание; Дедлайн задачи в формате дд.мм.гггг\n\n" +
			"Чтобы получить список своих задач, введи /todos\n\n" +
			"Чтобы удалить задачу введи /delete <ID задачи>",
	)
}

func (bot *TodoBot) CreateTodoHandler(ctx telebot.Context) error {
	taskArgs := ctx.Args()
	taskArgs = task.ParseTask(&taskArgs)

	check := task.CheckTask(taskArgs)
	if !check {
		return ctx.Send("Неверный формат!\n" +
			"Введите задачу в формате: /add Заголовок; Описание; Дедлайн задачи (дд.мм.гггг)")
	}

	existUser, err := bot.Users.FindOne(ctx.Chat().ID)
	if err != nil {
		log.Printf("Ошибка получения пользователя %v", err)
	}

	newTask := models.Task{
		Title:       taskArgs[0],
		Description: taskArgs[1],
		EndDate:     taskArgs[2],
		UserID:      existUser.ID,
	}

	err = bot.Tasks.Create(newTask)

	if err != nil {
		log.Printf("Ошибка создания задачи %v", err)
	}

	resStr := fmt.Sprintf("Создана задача %s: %s, до %s", newTask.Title, newTask.Description, newTask.EndDate)
	return ctx.Send(resStr)
}

func (bot *TodoBot) GetAllTodosHandler(ctx telebot.Context) error {
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
			`ID: %d
					Задача: %s
					Описание: %s
					Дедлайн: %s`,
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

func (bot *TodoBot) DeleteTodoHandler(ctx telebot.Context) error {
	taskArg := ctx.Args()

	if len(taskArg) == 0 {
		return ctx.Send("Вы не ввели ID задачи")
	}

	if len(taskArg) > 1 {
		return ctx.Send("Неверные данные, введите ID задачи" +
			"Чтобы посмотреть ID нужной задачи, введите /todos")
	}

	//existUser, err := bot.Users.FindOne(ctx.Chat().ID)
	//if err != nil {
	//	log.Printf("Ошибка получения пользователя %v", err)
	//}

	taskID, err := strconv.Atoi(taskArg[0])
	if err != nil {
		log.Printf("Ошибка удаления задачи. Неверный ID %v", err)
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
