package delivery

import (
	"fmt"
	"strconv"
	"time"

	"bot-tasker/internal/models"
	"bot-tasker/internal/service"

	"gopkg.in/telebot.v3"
)

type Handler struct {
	service  *service.Service
	msg      chan string
	commands []string
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		service:  s,
		msg:      make(chan string, 1),
		commands: []string{"/start", "/addTask", "/tasks", "/deleteTask"},
	}
}

func (h *Handler) StartHandler(c telebot.Context) error {
	user := models.User{
		Name:       c.Sender().Username,
		TelegramId: c.Sender().ID,
		FirstName:  c.Sender().FirstName,
		LastName:   c.Sender().LastName,
		ChatId:     c.Chat().ID,
	}
	if err := h.service.User.AddNewUser(user); err != nil {
		return err
	}
	return c.Send(fmt.Sprintf("Привет, %s. Список всех команд: /start, /addTask, /tasks, /deleteTask", user.FirstName))
}

func (h *Handler) AddTaskHandler(c telebot.Context) error {
	user, err := h.service.User.Get(c.Sender().ID)
	if err != nil {
		return err
	}
	if user == nil {
		return c.Send("Для начала нужно /start")
	}
	var task models.Task
	if err := c.Send("Напиши имя задачи."); err != nil {
		return err
	}
	task.Title = h.getUserInput()
	if err := c.Send("Напиши описание задачи."); err != nil {
		return err
	}
	task.Description = h.getUserInput()
	if err := c.Send("Укажи deadline (в днях)."); err != nil {
		return err
	}
	endDate, err := strconv.Atoi(h.getUserInput())
	if err != nil {
		return err
	}
	task.EndDate = time.Now().Add(time.Hour * 24 * time.Duration(endDate))
	task.UserID = user.ID
	if err := h.service.Task.CreateTask(task); err != nil {
		return err
	}
	return c.Send("Задача создана. Можешь посмотреть список всех задач при помощи команды /tasks")
}

func (h *Handler) getUserInput() string {
main:
	for {
		lastMsg := <-h.msg
		for _, comand := range h.commands {
			if lastMsg == comand {
				continue main
			}
		}
		return lastMsg
	}
}

func (h *Handler) AllTasksHandler(c telebot.Context) error {
	user, err := h.service.User.Get(c.Sender().ID)
	if err != nil {
		return err
	}
	if user == nil {
		return c.Send("Для начала нужно /start")
	}
	tasks, err := h.service.Task.GetAll(user.ID)
	if err != nil {
		return err
	}
	for _, task := range tasks {
		text := fmt.Sprintf("Задача #%d\nНазвание: %s\nОписание: %s\nДата создание: %s\nДата завершения: %s\n\n", task.ID, task.Title, task.Description, task.CreatedAt.Format("2006/01/02 15:04:05"), task.EndDate.Format("2006/01/02 15:04:05"))
		if err := c.Send(text); err != nil {
			return nil
		}
	}
	return nil
}

func (h *Handler) DeleteTaskHandler(c telebot.Context) error {
	args := c.Args()
	if len(args) == 0 {
		return c.Send("Пожалуйста укажите номер задачи: /deleteTask {id}")
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	user, err := h.service.User.Get(c.Sender().ID)
	if err != nil {
		return err
	}
	if user == nil {
		return c.Send("Для начала нужно /start")
	}
	if err := h.service.Task.Delete(user, uint(id)); err != nil {
		return err
	}
	return c.Send("Задача удалена.")
}

func (h *Handler) NewMsgHandler(c telebot.Context) error {
	h.msg <- c.Text()
	return nil
}
