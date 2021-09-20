package bot

import (
	"log"
	"os"
	"sync"

	"github.com/alexglazkov9/survgram/misc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var once sync.Once

// Database - Database Singleton that holds database reference
type Bot struct {
	bot     *tgbotapi.BotAPI
	updates map[string][]*tgbotapi.Update
}

var instance *Bot

// GetInstance - returns instance of database
func GetInstance() *Bot {
	once.Do(func() {
		instance = &Bot{}
		//Telegram bot setup
		bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN_DEV"))
		if err != nil {
			log.Panic(err)
		}

		instance.bot = bot
		instance.updates = make(map[string][]*tgbotapi.Update)
	})

	return instance
}

func (b *Bot) GetBot() *tgbotapi.BotAPI {
	return b.bot
}

func (b *Bot) HandleUpdate(update tgbotapi.Update) {

	if update.CallbackQuery != nil {
		cbData := misc.CallbackData{}
		cbData.FromJSON(update.CallbackQuery.Data)
		b.updates[cbData.Action] = append(b.updates[cbData.Action], &update)
	}

	//Possbile seprate commands handling to make commands case insensetive
	if update.Message != nil && len(update.Message.Text) > 0 {
		b.updates[update.Message.Text] = append(b.updates[update.Message.Text], &update)
	}
}

/* Pops next available update for one of the keywords.
Returns nil if there are no updates foir any of the keywords */
func (b *Bot) PopUpdate(keywords ...string) *tgbotapi.Update {
	var result *tgbotapi.Update
	result = nil
	for _, keyword := range keywords {
		q, ok := b.updates[keyword]
		if ok && len(q) > 0 {
			result = q[0]
			q = q[1:]
			b.updates[keyword] = q
			return result
		}
	}
	return result
}
