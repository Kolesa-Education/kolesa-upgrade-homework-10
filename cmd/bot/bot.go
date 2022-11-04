package bot

import (
	"log"
	"time"
	"upgrade/internal/models"
	"fmt"
	"strconv"

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
	//ctx.Send("Введите имя задачи:")
	
	taskTitle := ctx.Reply("Введите имя задачи:")
	fmt.Println(taskTitle)

	//if len(taskTitle) == 0 {
	//	return ctx.Send("Вы не ввели ваш вариант")
	//} else {
	//	ctx.Send(taskTitle)
	//}

	//ctx.Send("Введите описание задачи:")
	//taskDesc := ctx.Args()
	//
	//ctx.Send("Введите дедлайн задачи:")
	//taskEndDate := ctx.Args()
	//
	//ctx.Send(taskTitle)
	//ctx.Send(taskDesc)
	//ctx.Send(taskEndDate)
	//newTask := models.User{
	//	Title:       ctx.Sender().Username,
	//	Description: ctx.Chat().ID,
	//	UserId:  ctx.Sender().FirstName,
	//	EndDate:   ctx.Sender().LastName,
	//	//ChatId:     ctx.Chat().ID,
	//}
	//
	//err := bot.Users.Create(newTask)
	//
	//if err != nil {
	//	log.Printf("Ошибка создания пользователя %v", err)
	//}

	return ctx.Send("Test")
}

func (bot *UpgradeBot) GetUserTasksHandler(ctx telebot.Context) error {
	ctx.Send("Cписок всех Ваших задач: ")

	user, _ := bot.Users.FindOne(ctx.Chat().ID)

	result, err := bot.Tasks.GetAllUserTasks(int(user.ID))

	if err != nil {
		ctx.Send("Произошла ошибка")
		log.Fatalf("Ошибка: %v", err)
	}

	fmt.Println("Console log:")

	for _, record := range *result {
		recordId := strconv.Itoa(int(record.ID))
		UserId := strconv.Itoa(int(record.UserId))

		message := "Id задачи: " + recordId +
			"\r\nНазвание задачи: " + record.Title +
			"\r\nДедлайн: " + record.EndDate.String() +
			"\r\nUserID: " + UserId

		fmt.Println(message)
		ctx.Send(message)
	}

	return ctx.Send("Введите /addTask чтобы добавить новую задачу или /deleteTask {id} чтобы удалить существующую")
}

func (bot *UpgradeBot) DeleteTaskHandler(ctx telebot.Context) error {
	arguments := ctx.Args()

	if len(arguments) == 0 {
		return ctx.Send("Вы не ввели ваш id")
	}

	taskId, err := strconv.Atoi(arguments[0])

	if err != nil {
        return ctx.Send("Вы ввели не число")
    }

	user, _ := bot.Users.FindOne(ctx.Chat().ID)

	existTask, err := bot.Tasks.FindOneForUser(taskId, int(user.ID))

	if err != nil {
		return ctx.Send("Задачи с таким id не существует")
	}

	deleteErr := bot.Tasks.DeleteUserTask(*existTask, int(user.ID))

	if deleteErr != nil {
			log.Printf("Ошибка удаления задачи %v", err)
			return ctx.Send("Ошибка удаления задачи")
		}

		return ctx.Send("Задача удалена")
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