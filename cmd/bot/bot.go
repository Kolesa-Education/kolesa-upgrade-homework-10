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
}

var gameItems = [3]string{
    "камень",
    "ножницы",
    "бумага",
}

var winSticker = &telebot.Sticker{
    File: telebot.File{
        FileID: "CAACAgIAAxkBAAEGMEZjVspD4JulorxoH7nIwco5PGoCsAACJwADr8ZRGpVmnh4Ye-0RKgQ",
    },
    Width:    512,
    Height:   512,
    Animated: true,
}

var loseSticker = &telebot.Sticker{
    File: telebot.File{
        FileID: "CAACAgIAAxkBAAEGMEhjVsqoRriJRO_d-hrqguHNlLyLvQACogADFkJrCuweM-Hw5ackKgQ",
    },
    Width:    512,
    Height:   512,
    Animated: true,
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
    endDate, err := time.Parse("2006-01-02", "2018-01-20")
    if err != nil {
        log.Printf("Ошибка в дате %v", err)
    }
    newTask := models.Task{
        Name:       "Aaaaa",
        EndDate:    endDate,
        UserID:     335271283,
    }
    
    existUser, err := bot.Users.FindOne(335271283)
    

    if err != nil {
        log.Printf("Ошибка получения пользователя %v", err)
    }

    if existUser != nil {
        err := bot.Tasks.Create(newTask)

        if err != nil {
            log.Printf("Ошибка создания пользователя %v", err)
        }
    }

    return ctx.Send("Done")
}


func (bot *UpgradeBot) GetTasksHandler(ctx telebot.Context) error {
    existUser, err := bot.Users.FindOne(335271283)
    if err != nil {
        log.Printf("Ошибка получения пользователя %v", err)
    }

    if existUser != nil {
        tasks, err := bot.Tasks.FindByUserId(335271283)
        if err != nil {
            log.Printf("Ошибка при выводе задач %v", err)
        }
        log.Printf(strconv.Itoa(len(*tasks)))
        for i, task := range *tasks {
            log.Printf(strconv.Itoa(i+1) + " " + task.Name)
        }

    }

    return ctx.Send("Done")
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
