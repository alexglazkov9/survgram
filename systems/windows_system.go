package systems

import (
	"log"
	"strconv"

	"github.com/alexglazkov9/survgram/bot"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/interfaces"
	"github.com/alexglazkov9/survgram/items"
	"github.com/alexglazkov9/survgram/misc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

/* Handles InlineKeyboard menus and so called windows, e.g. Inventory, Character, etc */
type WindowsSystem struct {
	manager         *entity.Manager
	characterHelper interfaces.CharacterHelper
}

func NewWindowsSystem(manager *entity.Manager, characterHelper interfaces.CharacterHelper) *WindowsSystem {
	sys := &WindowsSystem{manager: manager, characterHelper: characterHelper}

	return sys
}

func (ws *WindowsSystem) Update(dt float64) {
	//Handle WINDOW_CLOSE
	for {
		u := bot.GetInstance().PopUpdate(
			misc.WINDOW_CLOSE,
		)
		if u == nil {
			break
		}
		bot.GetInstance().GetBot().DeleteMessage(tgbotapi.NewDeleteMessage(u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.MessageID))
	}
	//Handle INVENTORY window
	for {
		u := bot.GetInstance().PopUpdate(
			misc.INVENTORY_ITEM_SELECTED,
			misc.INVENTORY_DROP_ITEM,
			misc.INVENTORY_ITEM_USE,
		)
		if u == nil {
			break
		}

		chrctr := ws.characterHelper.GetCharacter(u.CallbackQuery.From.ID)
		if chrctr == nil || !chrctr.HasComponent("InventoryWindowComponent") {
			//Delete msg from the chat
			bot.GetInstance().GetBot().AnswerCallbackQuery(tgbotapi.NewCallback(u.CallbackQuery.ID, "This message is expired"))
			bot.GetInstance().GetBot().DeleteMessage(tgbotapi.NewDeleteMessage(u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.MessageID))
			continue
		}

		cbData := misc.CallbackData{}
		cbData.FromJSON(u.CallbackQuery.Data)

		inventoryWindow_C := chrctr.GetComponent("InventoryWindowComponent").(*components.InventoryWindowComponent)

		switch cbData.Action {
		case misc.INVENTORY_ITEM_SELECTED:
			select_item_id, _ := strconv.Atoi(cbData.Payload)
			inventoryWindow_C.SelectedItemID = &select_item_id
			inventoryWindow_C.SendUpdate()
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
			inventoryWindow_C.SelectedItemID = nil
			inventoryWindow_C.SendUpdate()
		case misc.INVENTORY_ITEM_USE:
			itm_id, _ := strconv.Atoi(cbData.Payload)
			if components.ItemID(itm_id).GetItem().GetType() == items.EQUIPMENT || components.ItemID(itm_id).GetItem().GetType() == items.WEAPON {
				inventory_C := chrctr.GetComponent("InventoryComponent").(*components.InventoryComponent)
				bundle, ok := inventory_C.TryFind(itm_id)
				log.Println("fetched from inventopry")
				if ok {
					equipment_C := chrctr.GetComponent("PlayerEquipmentComponent").(*components.PlayerEquipmentComponent)
					log.Println("equiping item")
					ok := equipment_C.TryEquip(components.ItemID(bundle.ID))

					if !ok {
						bot.GetInstance().GetBot().AnswerCallbackQuery(tgbotapi.NewCallback(u.CallbackQuery.ID, "⚠Could not equip item!"))
					} else {
						bundle.Qty--
						if bundle.Qty <= 0 {
							inventory_C.RemoveItem(bundle)
						}
						log.Println("equiping item ok")
						inventoryWindow_C.SendUpdate()
					}
				}
			}
		}
	}

	//Handle CHARACTER window
	for {
		u := bot.GetInstance().PopUpdate(
			misc.CHARACTER_SELECT_TAB,
			misc.CHARACTER_EQUIPMENT_SELECTED,
		)
		if u == nil {
			break
		}

		chrctr := ws.characterHelper.GetCharacter(u.CallbackQuery.From.ID)
		if chrctr == nil || !chrctr.HasComponent("CharacterWindowComponent") {
			//Delete msg from the chat
			bot.GetInstance().GetBot().AnswerCallbackQuery(tgbotapi.NewCallback(u.CallbackQuery.ID, "This message is expired"))
			bot.GetInstance().GetBot().DeleteMessage(tgbotapi.NewDeleteMessage(u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.MessageID))
			continue
		}

		cbData := misc.CallbackData{}
		cbData.FromJSON(u.CallbackQuery.Data)

		characterWindow_C := chrctr.GetComponent("CharacterWindowComponent").(*components.CharacterWindowComponent)

		switch cbData.Action {
		case misc.CHARACTER_SELECT_TAB:

			tab, _ := strconv.Atoi(cbData.Payload)
			if characterWindow_C.CurrentTab != components.CharacterWindowTab(tab) {
				characterWindow_C.CurrentTab = components.CharacterWindowTab(tab)
				characterWindow_C.SendUpdate()
			}
		case misc.CHARACTER_EQUIPMENT_SELECTED:
			item_id, _ := strconv.Atoi(cbData.Payload)
			characterWindow_C.SelectedItem = (*components.ItemID)(&item_id)
		}
	}

	//Handle InventoryWindowComponent update
	for _, e := range ws.manager.QueryEntities("InventoryWindowComponent") {
		inventoryWindow_C := e.GetComponent("InventoryWindowComponent").(*components.InventoryWindowComponent)
		if !inventoryWindow_C.IsSent {
			inventoryWindow_C.SendUpdate()
			inventoryWindow_C.IsSent = true
		}
	}
	//Handle CharacterWindowComponent update
	for _, e := range ws.manager.QueryEntities("CharacterWindowComponent") {
		characterWindow_C := e.GetComponent("CharacterWindowComponent").(*components.CharacterWindowComponent)
		if !characterWindow_C.IsSent {
			characterWindow_C.SendUpdate()
			characterWindow_C.IsSent = true
		}
	}
}
