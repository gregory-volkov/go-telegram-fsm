package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gregory-volkov/go-telegram-fsm"
	"github.com/gregory-volkov/go-telegram-fsm/example/handlers"
	"log"
	"os"
	"time"
)

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("Please set TELEGRAM_TOKEN environment variable.")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	store := fsm.NewInMemoryStore()
	machine := fsm.NewFSM(bot, store, "start", 5*time.Minute)

	machine.State("start").
		OnRegex(`(?i)^/start$`, handlers.HandleStart, "await_name").
		OnMessage("Hi", handlers.HandleHi, "greet").
		OnMessage("Hello", handlers.HandleHello, "greet")

	machine.State("greet").
		OnInlineButton("show_hello", handlers.HandleButton, "start").
		OnReplyButton("Option 1", handlers.HandleReplyButton, "start").
		OnReplyButton("Option 2", handlers.HandleReplyButton, "start")

	machine.State("wait_name").OnRegex(`.*`, handlers.HandleName, "start")

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if err := machine.ProcessUpdate(update); err != nil {
			log.Println("Error processing update:", err)
		}
	}
}
