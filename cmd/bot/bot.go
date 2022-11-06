package bot

import (
	"fmt"
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"log"
	"strings"
	"time"
	"upgrade/internal/models"

	"gopkg.in/telebot.v3"
)

type TodoBot struct {
	Bot   *telebot.Bot
	Users *models.UserModel
	Tasks *models.TaskModel
}

//var winSticker = &telebot.Sticker{
//	File: telebot.File{
//		FileID: "CAACAgIAAxkBAAEGMEZjVspD4JulorxoH7nIwco5PGoCsAACJwADr8ZRGpVmnh4Ye-0RKgQ",
//	},
//	Width:    512,
//	Height:   512,
//	Animated: true,
//}
//
//var loseSticker = &telebot.Sticker{
//	File: telebot.File{
//		FileID: "CAACAgIAAxkBAAEGUfFjZnDTpRmHFgABq_zncY60TpKZhlUAAgsBAAIWfGgDjgzxiPsp7OIrBA",
//	},
//	Width:    512,
//	Height:   512,
//	Animated: true,
//}

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

func (bot *TodoBot) CreateTodoHandler(ctx telebot.Context) error {
	taskArgs := ctx.Args()
	taskArgs = parseTask(&taskArgs)

	check := checkTask(taskArgs)
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

func checkTask(taskArgs []string) bool {
	if len(taskArgs) == 0 {
		log.Println("Нет аргументов!")
		return false
	}

	if len(taskArgs) < 3 {
		log.Println("Слишком мало аргументов")
		return false
	}

	regexpDate := pcre.MustCompile(`^(?:(?:31(\/|-|\.)(?:0?[13578]|1[02]))\1|(?:(?:29|30)(\/|-|\.)(?:0?[13-9]|1[0-2])\2))(?:(?:1[6-9]|[2-9]\d)?\d{2})$|^(?:29(\/|-|\.)0?2\3(?:(?:(?:1[6-9]|[2-9]\d)?(?:0[48]|[2468][048]|[13579][26])|(?:(?:16|[2468][048]|[3579][26])00))))$|^(?:0?[1-9]|1\d|2[0-8])(\/|-|\.)(?:(?:0?[1-9])|(?:1[0-2]))\4(?:(?:1[6-9]|[2-9]\d)?\d{2})$`, 0)
	match := regexpDate.MatcherString(taskArgs[2], 0).Matches()

	if !match {
		log.Println("Не прошла проверка по регулярке")
		return false
	}

	return true
}

func parseTask(args *[]string) []string {
	refString := strings.Join(*args, "")

	return strings.Split(refString, ";")
}

func (bot *TodoBot) HelpHandler(ctx telebot.Context) error {
	return ctx.Send(
		"**Список комманд**\n\n" +
			"Чтобы добавить задачу, введи её в формате\n" +
			"/add Заголовок; Описание; Дедлайн задачи в формате дд.мм.гггг\n\n" +
			"Чтобы получить список своих задач, введи /todos\n\n" +
			"Чтобы удалить задачу введи /delete <id задачи>",
	)
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
