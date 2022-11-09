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

	return ctx.Send("Привет, " + ctx.Sender().FirstName + "\nДля добавления новой задачи используй шаблон:\n/addTask домашка|на 100 баллов|15.08.2021\n/showTasks для просмотра всех задач\n/deleteTask 5 для удаления (5 - id задачи)")
}

func (bot *UpgradeBot) AddTaskHandler(ctx telebot.Context) error {

	allText := ctx.Text()

	if allText == "/addTask" {
		return ctx.Send("Для добавления новой задачи используй шаблон:\n/addTask домашка|на 100 баллов|15.08.2021")
	}

	clearText := strings.Replace(allText, "/addTask ", "", -1)

	vals := strings.Split(clearText, "|")

	if len(vals) != 3 {
		return ctx.Send("Введен некорректный шаблон")
	}

	date, err := time.Parse("02.01.2006", vals[2])

	if err != nil {
		return ctx.Send("Введена направильная дата")
	}

	newTask := models.Task{
		Title:       vals[0],
		Description: vals[1],
		UserId:      ctx.Sender().ID,
		EndDate:     date,
	}

	if err := bot.Tasks.Create(newTask); err != nil {
		log.Printf("Ошибка создания задачи %v", err)
		return ctx.Send("Ошибка создания задачи")
	}
	return ctx.Send("Задача добавлена")
}

func (bot *UpgradeBot) DeleteTaskHandler(ctx telebot.Context) error {

	if ctx.Text() == "/deleteTask" {
		return ctx.Send("Для удаления задачи используй шаблон\n/deleteTask 5 (5 - id задачи)")
	}

	args := ctx.Args()
	deleteId, err := strconv.ParseInt(args[0], 0, 64)

	if err != nil {
		return ctx.Send("Введен направильный id")
	}

	if len(args) > 1 {
		return ctx.Send("Введен направильный id")
	}

	if err := bot.Tasks.Delete(deleteId, ctx.Sender().ID); err != nil {
		log.Fatalf("Ошибка выполнения запроса пользователя %v", err)
	}

	return ctx.Send("Задача удалена")

}

func (bot *UpgradeBot) ShowTasksHandler(ctx telebot.Context) error {

	tasks, err := bot.Tasks.FindAll(ctx.Sender().ID)

	if err != nil {
		log.Fatalf("Ошибка загрузки %v", err)
	}

	var (
		Tasks     []string
		task_info string
	)

	for _, el := range tasks {
		task_info = fmt.Sprintf(`%d. %s - %s, успеть до: %s`, el.ID, el.Title, el.Description, el.EndDate.Format("02.01.2006"))
		Tasks = append(Tasks, task_info)
	}

	resString := strings.Join(Tasks, "\n")

	return ctx.Send(resString)
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
