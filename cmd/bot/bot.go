package bot

import (
    "log"
    "time"
    "upgrade/internal/models"
    "strconv"
    "gopkg.in/telebot.v3"
)


type UpgradeBot struct {
    Bot   *telebot.Bot
    Users *models.UserModel
    Tasks *models.TaskModel
    AddTaskState models.AddTaskState
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
    bot.AddTaskState.CurrentState = models.Name
    return ctx.Send("Введите имя задачи...")
}


func (bot *UpgradeBot) DeleteTaskHandler(ctx telebot.Context) error {
    args := ctx.Args()

    if len(args) == 0 {
        return ctx.Send("Ошибка. Нужно ввести ID задания...")
    }
    for _, id := range args{
        err := bot.Tasks.Delete(id, models.Standart)
        if err != nil{
            return ctx.Send("Ошибка с ID:" + id)
        }
    }
    return ctx.Send("Done")
}

func (bot *UpgradeBot) GetTasksHandler(ctx telebot.Context) error {
    existUser, err := bot.Users.FindOne(ctx.Chat().ID)
    result := ""
    if err != nil {
        log.Printf("Ошибка получения пользователя %v", err)
    }

    if existUser != nil {
        tasks, err := bot.Tasks.FindByUserId(ctx.Chat().ID)
        if err != nil {
            return ctx.Send("Возникла ошибка при выводе задач...")
        }
        log.Printf(strconv.Itoa(len(*tasks)))
        if len(*tasks) == 0{
            return ctx.Send("No Tasks...")
        }
        for i, task := range *tasks {
            formatedDate := task.EndDate.Format("01/02/2006")
            result += strconv.Itoa(i+1) + ". " + task.Name + " | " + formatedDate+ "\n"
        }

    }

    return ctx.Send(result)
}


func (bot *UpgradeBot) GeneralHandler(ctx telebot.Context) error {
    var result error
    addTaskResult := bot.AddTaskState.HandleState(ctx)
    switch addTaskResult {
    case "Action::Add":
        endDate, _ := time.Parse("01/02/2006", bot.AddTaskState.Storage["endDate"])
        newTask := models.Task{
            Name:       bot.AddTaskState.Storage["name"],
            EndDate:    endDate,
            UserID:     uint(ctx.Chat().ID),
        }
        return ctx.Send(addTask(newTask, bot, ctx))
    case "":
        result = nil
    default:
        return ctx.Send(addTaskResult)
    }

        


    return result
}

func addTask(task models.Task, bot *UpgradeBot, ctx telebot.Context)string{
    existUser, err := bot.Users.FindOne(ctx.Chat().ID)
    if err != nil {
        return "Ошибка получения пользователя"
    }
    if existUser != nil {
        err := bot.Tasks.Create(task)

        if err != nil {
            return "Ошибка создания пользователя"
        }
        return "Успешно!"
    }
    return "Пользователя нет в базе"
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
