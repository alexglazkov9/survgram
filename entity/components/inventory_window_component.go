package components

import (
	"fmt"

	"github.com/alexglazkov9/survgram/bot"
	"github.com/alexglazkov9/survgram/items"
	"github.com/alexglazkov9/survgram/misc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type InventoryWindowComponent struct {
	BaseComponent

	SelectedItemID *int
	IsSent         bool
	Message        *tgbotapi.Message
}

func (iwc *InventoryWindowComponent) SendUpdate() {
	player_C := iwc.Parent.GetComponent("PlayerComponent").(*PlayerComponent)
	inventory_C := iwc.Parent.GetComponent("InventoryComponent").(*InventoryComponent)

	text := "Inventory"
	tgkb := &misc.TGInlineKeyboard{Columns: 2, IsClosable: true}

	if iwc.SelectedItemID != nil {
		cb_data := misc.CallbackData{Action: misc.INVENTORY_DROP_ITEM, Payload: fmt.Sprint(*iwc.SelectedItemID)}
		tgkb.AddHeaderButton("🚮Drop Item", cb_data.JSON())
	}

	for _, item_bundle := range inventory_C.GetItems() {
		if iwc.SelectedItemID != nil && item_bundle.ID == *iwc.SelectedItemID {
			//Add EQUIP button if Weapon or Equipment
			if item_bundle.GetItem().GetType() == items.EQUIPMENT || item_bundle.GetItem().GetType() == items.WEAPON {
				cb_data := misc.CallbackData{Action: misc.INVENTORY_ITEM_USE, Payload: fmt.Sprint(item_bundle.ID)}
				tgkb.AddHeaderButton("🔼Equip", cb_data.JSON())
			}
			text = item_bundle.GetItem().GetFormattedItem(true)
			cb_data := misc.CallbackData{Action: misc.INVENTORY_ITEM_SELECTED, Payload: fmt.Sprint(item_bundle.ID)}
			tgkb.AddButton(
				fmt.Sprintf("👁%s (%d)", item_bundle.GetItem().GetName(), item_bundle.Qty),
				cb_data.JSON(),
			)
		} else {
			cb_data := misc.CallbackData{Action: misc.INVENTORY_ITEM_SELECTED, Payload: fmt.Sprint(item_bundle.ID)}
			tgkb.AddButton(
				fmt.Sprintf("%s (%d)", item_bundle.GetItem().GetName(), item_bundle.Qty),
				cb_data.JSON(),
			)
		}
	}
	for i := len(inventory_C.Items); i < inventory_C.Slots; i++ {
		cb_data := misc.CallbackData{Action: "", Payload: ""}
		tgkb.AddButton("-", cb_data.JSON())
	}

	var msg tgbotapi.Chattable
	if iwc.Message != nil {
		msg_t := tgbotapi.NewEditMessageText(player_C.ChatID, iwc.Message.MessageID, text)
		msg_t.ReplyMarkup = tgkb.Generate()
		msg_t.ParseMode = "markdown"
		msg = msg_t
	} else {
		msg_t := tgbotapi.NewMessage(player_C.ChatID, text)
		msg_t.ReplyMarkup = tgkb.Generate()
		msg_t.ParseMode = "markdown"
		msg = msg_t
	}

	sent_msg, _ := bot.GetInstance().GetBot().Send(msg)
	iwc.Message = &sent_msg
}
