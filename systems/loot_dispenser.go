package systems

import (
	"strconv"

	"github.com/alexglazkov9/survgram/bot"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/interfaces"
	"github.com/alexglazkov9/survgram/misc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type LootDispenserSystem struct {
	manager         *entity.Manager
	characterHelper interfaces.CharacterHelper
}

func NewLootDispenserSystem(manager *entity.Manager, characterHelper interfaces.CharacterHelper) *LootDispenserSystem {
	sys := &LootDispenserSystem{manager: manager, characterHelper: characterHelper}

	return sys
}

func (lds *LootDispenserSystem) Update(dt float64) {
	for {
		u := bot.GetInstance().PopUpdate(
			misc.DISMISS_LOOT,
			misc.SELECT_LOOT_ITEM,
			misc.PICK_UP_ITEM,
			misc.PICK_UP_ALL_ITEMS,
		)
		if u == nil {
			break
		}

		cbData := misc.CallbackData{}
		cbData.FromJSON(u.CallbackQuery.Data)
		e_id, _ := strconv.Atoi(cbData.ID)

		e := lds.manager.GetEntityByID(e_id)
		lootDispenser_C, lootDispenser_C_ok := e.GetComponent("LootDispenserComponent").(*components.LootDispenserComponent)

		if lootDispenser_C_ok {
			switch cbData.Action {
			//DISMISS LOOT
			case misc.DISMISS_LOOT:
				lootDispenser_C.Expire()
			//ITEM SELECTED
			case misc.SELECT_LOOT_ITEM:
				for _, itm := range lootDispenser_C.LootDrop {
					itm_Id, _ := strconv.Atoi(cbData.Payload)
					if itm.ID == itm_Id {
						lootDispenser_C.SelectedItem = &itm
						lootDispenser_C.UpdateMessage()
					}
				}
			//PICK UP ITEM
			case misc.PICK_UP_ITEM:
				chrctr := lds.characterHelper.GetCharacter(u.CallbackQuery.From.ID)
				inventory_C := chrctr.GetComponent("InventoryComponent").(*components.InventoryComponent)
				inventory_C.AddItems(*lootDispenser_C.SelectedItem)
				for i, item := range lootDispenser_C.LootDrop {
					if item.ID == lootDispenser_C.SelectedItem.ID {
						lootDispenser_C.LootDrop[i] = lootDispenser_C.LootDrop[len(lootDispenser_C.LootDrop)-1]
						lootDispenser_C.LootDrop = lootDispenser_C.LootDrop[:len(lootDispenser_C.LootDrop)-1]
						lootDispenser_C.SelectedItem = nil
						lootDispenser_C.UpdateMessage()
					}
				}
			//PICK UP ALL ITEMS
			case misc.PICK_UP_ALL_ITEMS:
			}
		}
	}
	for _, p := range lds.manager.QueryEntities("LootDispenserComponent") {
		lootDispenser_C := p.GetComponent("LootDispenserComponent").(*components.LootDispenserComponent)
		switch lootDispenser_C.State {
		case components.ADDED:
			//Send Message
			text := "Loot"
			player_C := p.GetComponent("PlayerComponent").(*components.PlayerComponent)
			msg := tgbotapi.NewMessage(player_C.ChatID, text)
			msg.ReplyMarkup = lootDispenser_C.GenerateInlineKeyboard()
			message, _ := bot.GetInstance().GetBot().Send(msg)
			lootDispenser_C.Message = &message
			lootDispenser_C.State = components.PICKING
		}
	}
}
