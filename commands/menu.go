package commands

import (
	"github.com/alexglazkov9/survgram/game"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Menu(bot *tgbotapi.BotAPI, update tgbotapi.Update, game game.Game) {
	if update.Message.Chat.IsPrivate() {
		char := game.CharacterManager.GetCharacter(update.Message.From.ID)
		if char != nil {
			replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Go to...", "goto"),
				),
			)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Menu")
			msg.ReplyMarkup = replyMarkup
			bot.Send(msg)
			return
		}

	} else {
		// Delete message if it is sent somwhere other than private chat
		bot.DeleteMessage(tgbotapi.DeleteMessageConfig{MessageID: update.Message.MessageID, ChatID: update.Message.Chat.ID})
	}
}
