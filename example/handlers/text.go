package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gregory-volkov/go-telegram-fsm/fsm"
)

func HandleStart(ctx *fsm.Context, u tgbotapi.Update) (fsm.StateID, error) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Welcome! What's your name?")
	_, err := ctx.Bot.Send(msg)
	return "wait_name", err
}

func HandleName(ctx *fsm.Context, u tgbotapi.Update) (fsm.StateID, error) {
	name := u.Message.Text
	ctx.Vars["username"] = name
	_ = ctx.Store.Set(ctx.Ctx, ctx.Key, ctx.Vars)
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Nice to meet you, "+name+"! Returning to the `start` state.")
	_, err := ctx.Bot.Send(msg)
	return "start", err
}

func HandleHi(ctx *fsm.Context, u tgbotapi.Update) (fsm.StateID, error) {
	chatID := u.Message.Chat.ID

	button := tgbotapi.NewInlineKeyboardButtonData("Say Hello 👋", "show_hello")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(button),
	)

	msg := tgbotapi.NewMessage(chatID, "Click the button below:")
	msg.ReplyMarkup = keyboard
	_, err := ctx.Bot.Send(msg)

	return "greet", err
}

func HandleHello(ctx *fsm.Context, u tgbotapi.Update) (fsm.StateID, error) {
	chatID := u.Message.Chat.ID

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Option 1"),
			tgbotapi.NewKeyboardButton("Option 2"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, "Choose an option:")
	msg.ReplyMarkup = keyboard
	_, err := ctx.Bot.Send(msg)

	return "greet", err
}
