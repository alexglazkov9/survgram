package main

import (
	"log"

	"github.com/alexglazkov9/survgram/bot"
	"github.com/alexglazkov9/survgram/game"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	_ = game.New(bot.GetInstance().GetBot())

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetInstance().GetBot().GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		bot.GetInstance().HandleUpdate(update)
	}
}
