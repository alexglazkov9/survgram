package systems

import (
	"fmt"

	"github.com/alexglazkov9/survgram/activities"
	"github.com/alexglazkov9/survgram/bot"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func SendMainMenuKeyboard(e *entity.Entity) {
	player_C := e.GetComponent("PlayerComponent").(*components.PlayerComponent)
	loc := activities.GetLocations().GetLocation(player_C.CurrentLocation)
	txt := fmt.Sprintf("You are at the %s", loc.Name)
	kb := tgbotapi.NewReplyKeyboard()
	kb.Keyboard = append(kb.Keyboard, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(bot.MENU_OPTION_EXPEDITION),
		tgbotapi.NewKeyboardButton(bot.MENU_OPTION_NPCS),
	))
	kb.Keyboard = append(kb.Keyboard, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(bot.MENU_OPTION_CHARACTER),
		tgbotapi.NewKeyboardButton(bot.MENU_OPTION_HIDEOUT),
	))
	kb.Keyboard = append(kb.Keyboard, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(bot.MAP),
	))
	msg := tgbotapi.NewMessage(player_C.ChatID, txt)
	msg.ReplyMarkup = kb
	bot.GetInstance().GetBot().Send(msg)
}
