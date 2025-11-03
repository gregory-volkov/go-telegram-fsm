package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gregory-volkov/go-telegram-fsm"
)

func HandleButton(ctx *fsm.Context, u tgbotapi.Update) (fsm.StateID, error) {
	cb := u.CallbackQuery
	chatID := cb.Message.Chat.ID

	answer := tgbotapi.NewCallback(cb.ID, "Button pressed!")
	if _, err := ctx.Bot.Request(answer); err != nil {
		return "", err
	}

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Hello back to you, %s! 🎉", cb.From.FirstName))
	_, err := ctx.Bot.Send(msg)

	return "start", err
}

func HandleReplyButton(ctx *fsm.Context, u tgbotapi.Update) (fsm.StateID, error) {
	chatID := u.Message.Chat.ID
	selection := u.Message.Text

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("You selected: %s", selection))
	_, err := ctx.Bot.Send(msg)

	return "start", err
}
