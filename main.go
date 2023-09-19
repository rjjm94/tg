//main.go// Path: main.go

package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	db "tg/db"
	"tg/handlers"
)

func main() {
	bot, updates, err := initializeBot()
	if err != nil {
		log.Fatal(err)
	}

	handleUpdates(bot, updates)
}

func initializeBot() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel, error) {
	bot, err := tgbotapi.NewBotAPI("6609686170:AAE-yBE_s3NmxUu_q5Ir62iQ-OaLur1BUCU")
	if err != nil {
		return nil, nil, err
	}

	bot.Debug = true // Change this to false in production

	db.Connect() // Connect to your MongoDB database

	handlers.SetBot(bot) // Set the bot in your handlers package

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return nil, nil, err
	}

	return bot, updates, nil
}

func handleUpdates(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil || update.CallbackQuery != nil {
			response := handlers.HandleMessage(&update)
			if response != nil {
				bot.Send(response)
			}
		}
	}
}
