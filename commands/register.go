package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/alexglazkov9/survgram/bot"
	"github.com/alexglazkov9/survgram/game"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	MINNAMELENGTH = 7
	MAXNAMELENGTH = 15
)

//Register - Verifies and creates a character
func Register(update tgbotapi.Update, game game.Game) {
	if update.Message.Chat.IsPrivate() {
		// Check if the player already has a character
		if game.CharacterManager.GetCharacter(update.Message.From.ID) != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You already have a character")
			bot.GetInstance().GetBot().Send(msg)
			return
		}

		// Fetch and validate name
		name := update.Message.CommandArguments()
		if isNameValid(name) {
			if game.CharacterManager.CreateCharacter(update.Message.From.ID, update.Message.Chat.ID, name) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Welcome to the game %s", name))
				bot.GetInstance().GetBot().Send(msg)
				return
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Something went wrong...")
				bot.GetInstance().GetBot().Send(msg)
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("The name \"%s\" is invalid", name))
			bot.GetInstance().GetBot().Send(msg)
			log.Printf("Invalid name: %s", name)
		}
	} else {
		// Delete message if it is sent somwhere other than private chat
		bot.GetInstance().GetBot().DeleteMessage(tgbotapi.DeleteMessageConfig{MessageID: update.Message.MessageID, ChatID: update.Message.Chat.ID})
	}
}

func isNameValid(name string) bool {
	if len(name) >= MINNAMELENGTH && !strings.ContainsAny(name, "  | - | + | ! | ? | . | ,") && len(name) <= MAXNAMELENGTH {
		return true
	}
	return false
}
