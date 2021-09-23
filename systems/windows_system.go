package systems

import (
	"fmt"
	"strconv"

	"github.com/alexglazkov9/survgram/bot"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/interfaces"
	"github.com/alexglazkov9/survgram/items"
	"github.com/alexglazkov9/survgram/misc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type WindowsSystem struct {
	manager         *entity.Manager
	characterHelper interfaces.CharacterHelper
}

func NewWindowsSystem(manager *entity.Manager, characterHelper interfaces.CharacterHelper) *WindowsSystem {
	sys := &WindowsSystem{manager: manager, characterHelper: characterHelper}

	return sys
}

func (ws *WindowsSystem) Update(dt float64) {
	for {
		u := bot.GetInstance().PopUpdate(
			misc.INVENTORY_ITEM_SELECTED,
			misc.INVENTORY_DROP_ITEM,
			misc.WINDOW_CLOSE,
		)
		if u == nil {
			break
		}
		chrctr := ws.characterHelper.GetCharacter(u.CallbackQuery.From.ID)
		if chrctr == nil {
			//TODO Handle missing character
			continue
		}
		cbData := misc.CallbackData{}
		cbData.FromJSON(u.CallbackQuery.Data)
		if cbData.Action == misc.WINDOW_CLOSE {
			bot.GetInstance().GetBot().DeleteMessage(tgbotapi.NewDeleteMessage(u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.MessageID))
		}

		switch cbData.Action {
		case misc.INVENTORY_ITEM_SELECTED:
			inventoryWindow_C := chrctr.GetComponent("InventoryWindowComponent").(*components.InventoryWindowComponent)
			select_item_id, _ := strconv.Atoi(cbData.Payload)
			inventoryWindow_C.SelectedItemID = &select_item_id
			bot.GetInstance().GetBot().Send(editInventoryMessage(chrctr, u))
		case misc.INVENTORY_DROP_ITEM:
			itm_id, _ := strconv.Atoi(cbData.Payload)
			inventory_C := chrctr.GetComponent("InventoryComponent").(*components.InventoryComponent)
			for i, itm_b := range inventory_C.GetItems() {
				if itm_b.ID == itm_id {
					inventory_C.Items = append(inventory_C.Items[:i], inventory_C.Items[i+1:]...)
					bot.GetInstance().GetBot().AnswerCallbackQuery(tgbotapi.NewCallback(u.CallbackQuery.ID, "You have dropped item"))
					break
				}
			}
			inventoryWindow_C := chrctr.GetComponent("InventoryWindowComponent").(*components.InventoryWindowComponent)
			inventoryWindow_C.SelectedItemID = nil
			bot.GetInstance().GetBot().Send(editInventoryMessage(chrctr, u))
		}
	}
	for _, e := range ws.manager.QueryEntities("InventoryWindowComponent") {
		inventoryWindow_C := e.GetComponent("InventoryWindowComponent").(*components.InventoryWindowComponent)
		if !inventoryWindow_C.IsSent {

			bot.GetInstance().GetBot().Send(generateInventoryMessage(e))
			inventoryWindow_C.IsSent = true
		}
	}
}

func editInventoryMessage(e *entity.Entity, u *tgbotapi.Update) tgbotapi.Chattable {
	player_C := e.GetComponent("PlayerComponent").(*components.PlayerComponent)
	inventory_C := e.GetComponent("InventoryComponent").(*components.InventoryComponent)
	inventoryWindow_C := e.GetComponent("InventoryWindowComponent").(*components.InventoryWindowComponent)

	var text string
	text = "Inventory"
	tgkb := &misc.TGInlineKeyboard{Columns: 2, IsClosable: true}

	if inventoryWindow_C.SelectedItemID != nil {
		cb_data := misc.CallbackData{Action: misc.INVENTORY_DROP_ITEM, Payload: fmt.Sprint(*inventoryWindow_C.SelectedItemID)}
		tgkb.AddHeaderButton("🚮Drop Item", cb_data.JSON())
	}

	for _, item_bundle := range inventory_C.GetItems() {
		if inventoryWindow_C.SelectedItemID != nil && item_bundle.ID == *inventoryWindow_C.SelectedItemID {
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

	msg := tgbotapi.NewEditMessageText(player_C.ChatID, u.CallbackQuery.Message.MessageID, text)
	msg.ReplyMarkup = tgkb.Generate()
	msg.ParseMode = "markdown"
	return msg
}

func generateInventoryMessage(e *entity.Entity) tgbotapi.Chattable {
	player_C := e.GetComponent("PlayerComponent").(*components.PlayerComponent)
	inventory_C := e.GetComponent("InventoryComponent").(*components.InventoryComponent)

	text := "Inventory"
	tgkb := &misc.TGInlineKeyboard{Columns: 2, IsClosable: true}

	for _, item_bundle := range inventory_C.GetItems() {
		cb_data := misc.CallbackData{Action: misc.INVENTORY_ITEM_SELECTED, Payload: fmt.Sprint(item_bundle.ID)}
		tgkb.AddButton(
			fmt.Sprintf("%s (%d)", item_bundle.GetItem().GetName(), item_bundle.Qty),
			cb_data.JSON(),
		)
	}
	for i := len(inventory_C.Items); i < inventory_C.Slots; i++ {
		cb_data := misc.CallbackData{Action: "", Payload: ""}
		tgkb.AddButton("-", cb_data.JSON())
	}

	msg := tgbotapi.NewMessage(player_C.ChatID, text)
	msg.ReplyMarkup = tgkb.Generate()
	msg.ParseMode = "markdown"
	return msg
}
