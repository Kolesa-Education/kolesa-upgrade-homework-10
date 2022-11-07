package models

import (
	"gopkg.in/telebot.v3"
	"time"
)


type State struct{
	Storage map[string]string
}


type AddTaskState struct{
	State
	CurrentState AddTaskCurrentState
}

func (s *AddTaskState) HandleState(ctx telebot.Context)string{
	text := ctx.Text()
	switch s.CurrentState {
    case Name:
        s.Storage["name"] = text
        s.CurrentState = EndDate
        return "Введите дату в формате ДД/ММ/ГГГГ..."
    case EndDate:
        _, err := time.Parse("01/02/2006", text)
        if err != nil{
            return "Неверный формат..."
        }
		
        s.Storage["endDate"] = text
		s.CurrentState = None
		return "Action::Add"
    }
	return ""
	
}





type AddTaskCurrentState int

const (
    None AddTaskCurrentState = iota
	Name
	EndDate
)