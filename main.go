package main

import (
	"log"
	"os"

	"github.com/alexglazkov9/survgram/commands"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	//Telegram bot setup
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN_DEV"))
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		log.Printf(update.Message.Text)

		if update.Message != nil && update.Message.IsCommand() {
			switch update.Message.Command() {
			case "reg": // Create character
				go commands.Register(bot, update)
			}
		}
	}
}
