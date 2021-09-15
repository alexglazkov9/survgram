package commands

import (
	"fmt"
	"log"

	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/game"
	"github.com/alexglazkov9/survgram/misc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Inventory(bot *tgbotapi.BotAPI, update tgbotapi.Update, game game.Game) {
	if update.Message.Chat.IsPrivate() {
		char := game.CharacterManager.GetCharacter(update.Message.From.ID)
		//TODO Figure out nil check for not registered users
		if char != nil {
			tgkb := &misc.TGInlineKeyboard{Columns: 2}
			inv_C := char.GetComponent("InventoryComponent").(*components.InventoryComponent)
			for _, item_bundle := range inv_C.GetItems() {
				cb_data := misc.CallbackData{Action: "", Payload: ""}
				tgkb.AddButton(
					fmt.Sprintf("%s (%d)", item_bundle.GetItem().GetName(), item_bundle.Qty),
					cb_data.JSON(),
				)
			}
			for i := len(inv_C.Items); i < inv_C.Slots; i++ {
				cb_data := misc.CallbackData{Action: "", Payload: ""}
				tgkb.AddButton("-", cb_data.JSON())
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Inventory")
			msg.ReplyMarkup = tgkb.Generate()
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
			return
		}

	} else {
		// Delete message if it is sent somwhere other than private chat
		bot.DeleteMessage(tgbotapi.DeleteMessageConfig{MessageID: update.Message.MessageID, ChatID: update.Message.Chat.ID})
	}
}
