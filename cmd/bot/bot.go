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

type NewTask struct {
	Title string
	Desc string
	EndDate time.Time
}

var task = NewTask{}

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
	ctx.Send("Введите имя задачи с помощью команды: /taskName имя_задачи")

	return nil
}

func (bot *UpgradeBot) AddTaskNameHandler(ctx telebot.Context) error {
	taskTitle := ctx.Args()

	if len(taskTitle) == 0 { 
		return ctx.Send("Вы не ввели ваш вариант")
	}

	task.Title = taskTitle[0]
	
	return ctx.Send("Введите описание задачи с помощью команды: /taskDesc описание_задачи")
}

func (bot *UpgradeBot) AddTaskDescriptionHandler(ctx telebot.Context) error {
	taskDesc := ctx.Args()

	if len(taskDesc) == 0 { 
		return ctx.Send("Вы не ввели ваш вариант")
	}
	
	task.Desc = taskDesc[0]

	return ctx.Send("Введите дату завершения задачи с помощью команды: /taskEndDate YYYY-MM-DD")
}

func (bot *UpgradeBot) AddTaskEndDateHandler(ctx telebot.Context) error {
	taskEndDate := ctx.Args()

	if len(taskEndDate) == 0 { 
		return ctx.Send("Вы не ввели ваш вариант")
	}
	
	date, error := time.Parse("2006-01-02", taskEndDate[0])

	if error != nil {
        panic(error)
    }

	task.EndDate = date

	return ctx.Send("Введите /endTaskCreation чтобы завершить добавление задачи")
}

func (bot *UpgradeBot) EndTaskCreationHandler(ctx telebot.Context) error {
	user, _ := bot.Users.FindOne(ctx.Chat().ID)

	if len(task.Title) == 0 || len(task.Desc) == 0 || task.EndDate.IsZero() {
		return ctx.Send("Не все необходимые данные заполнены")
	}

	newTask := models.Task{
		Title:       task.Title,
		Description: task.Desc,
		EndDate: task.EndDate,
		UserId: uint(user.ID),
	}

	err := bot.Tasks.Create(newTask)

	if err != nil {
		log.Printf("Ошибка создания пользователя %v", err)
	}

	return ctx.Send("Задача добавлена")
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