package systems

import (
	"fmt"
	"log"
	"strings"

	"github.com/alexglazkov9/survgram/activities"
	"github.com/alexglazkov9/survgram/bot"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/resources"
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
		tgbotapi.NewKeyboardButton(bot.MENU_OPTION_MAP),
	))
	msg := tgbotapi.NewMessage(player_C.ChatID, txt)
	msg.ReplyMarkup = kb
	bot.GetInstance().GetBot().Send(msg)
	menu_C := e.GetComponent("MenuComponent").(*components.MenuComponent)
	menu_C.Menus.Clear()
	menu_C.Menus.Push(msg)
}

func GetMainMenu(e *entity.Entity) tgbotapi.MessageConfig {
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
		tgbotapi.NewKeyboardButton(bot.MENU_OPTION_MAP),
	))
	msg := tgbotapi.NewMessage(player_C.ChatID, txt)
	msg.ReplyMarkup = kb
	return msg
}

func GetExpeditionMenu(e *entity.Entity) tgbotapi.MessageConfig {
	player_C := e.GetComponent("PlayerComponent").(*components.PlayerComponent)
	loc := activities.GetLocations().GetLocation(player_C.CurrentLocation)
	msgTxt := "Possible activities in this location:\n"
	for _, act := range loc.PossibleActivities {
		msgTxt += fmt.Sprintf("%.0f%% %s \n", act.ActivityChance*100, strings.Title(string(act.Type)))
	}
	kb := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(bot.EXPEDITION_START),
		tgbotapi.NewKeyboardButton(bot.EXPEDITION_ACTIVITIES_INFO),
		tgbotapi.NewKeyboardButton(bot.MENU_BACK),
	))
	msg := tgbotapi.NewMessage(player_C.ChatID, msgTxt)
	msg.ReplyMarkup = kb
	return msg
}

func GetHideoutMenu(e *entity.Entity) tgbotapi.Chattable {
	log.Println("Sending hideout")
	player_C := e.GetComponent("PlayerComponent").(*components.PlayerComponent)
	msg := tgbotapi.NewPhotoUpload(player_C.ChatID, "./resources/images/hideout.png")
	msg.Caption = resources.HIDEOUT_WELCOME_TEXT
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(bot.HIDEOUT_STORAGE),
			tgbotapi.NewKeyboardButton(bot.HIDEOUT_CRAFTING_TABLE),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(bot.MENU_BACK),
		),
	)
	msg.ParseMode = "markdown"
	msg.ReplyMarkup = kb
	return msg
}

func GetExpeditionQuickMenu(e *entity.Entity) tgbotapi.MessageConfig {
	player_C := e.GetComponent("PlayerComponent").(*components.PlayerComponent)
	loc := activities.GetLocations().GetLocation(player_C.CurrentLocation)
	msg := tgbotapi.NewMessage(player_C.ChatID, fmt.Sprintf("%s at %s", resources.EXPEDITION_ENTERED_TEXT, loc.Name))
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(bot.EXPEDITION_RUNNING_MENU_CHARACTER),
			tgbotapi.NewKeyboardButton(bot.EXPEDITION_RUNNING_MENU_INVENTORY),
			tgbotapi.NewKeyboardButton(bot.EXPEDITION_RUNNING_MENU_ABILITIES),
		),
	)
	msg.ReplyMarkup = kb
	return msg
}

func GetQuickMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(bot.EXPEDITION_RUNNING_MENU_CHARACTER),
			tgbotapi.NewKeyboardButton(bot.EXPEDITION_RUNNING_MENU_INVENTORY),
			tgbotapi.NewKeyboardButton(bot.EXPEDITION_RUNNING_MENU_ABILITIES),
		),
	)
	return kb
}

func StartExpedition(e *entity.Entity) {

}
