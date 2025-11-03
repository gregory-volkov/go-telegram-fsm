package fsm

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
)

type stateBuilder struct {
	f  *FSM
	id StateID
}

func (f *FSM) State(id StateID) *stateBuilder {
	return &stateBuilder{f: f, id: id}
}

func (s *stateBuilder) OnRegex(regex string, h Handler, to StateID) *stateBuilder {
	r := regexp.MustCompile(regex)
	pred := func(u tgbotapi.Update) bool {
		return u.Message != nil && r.MatchString(u.Message.Text)
	}
	tr := Transition{From: s.id, To: to, Run: h, If: pred}
	s.f.RegisterTransition(tr)
	return s
}

func (s *stateBuilder) OnMessage(message string, h Handler, to StateID) *stateBuilder {
	tr := Transition{From: s.id, To: to, Run: h, If: func(u tgbotapi.Update) bool { return u.Message.Text == message }}
	s.f.RegisterTransition(tr)
	return s
}

func (s *stateBuilder) OnInlineButton(data string, h Handler, to StateID) *stateBuilder {
	pred := func(u tgbotapi.Update) bool {
		return u.CallbackQuery != nil && u.CallbackQuery.Data == data
	}
	tr := Transition{From: s.id, To: to, Run: h, If: pred}
	s.f.RegisterTransition(tr)
	return s
}

func (s *stateBuilder) OnReplyButton(buttonText string, h Handler, to StateID) *stateBuilder {
	pred := func(u tgbotapi.Update) bool {
		return u.Message != nil && u.Message.Text == buttonText
	}
	tr := Transition{From: s.id, To: to, Run: h, If: pred}
	s.f.RegisterTransition(tr)
	return s
}
