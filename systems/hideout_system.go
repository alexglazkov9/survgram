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
	"github.com/alexglazkov9/survgram/resources"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type HideoutSystem struct {
	currentDt float64

	manager         *entity.Manager
	characterHelper interfaces.CharacterHelper
}

func NewHideoutSystem(manager *entity.Manager, characterHelper interfaces.CharacterHelper) *HideoutSystem {
	sys := &HideoutSystem{manager: manager, characterHelper: characterHelper}

	return sys
}

func (hs *HideoutSystem) Update(dt float64) {
	for {
		u := bot.GetInstance().PopUpdate(bot.HIDEOUT_CRAFTING_TABLE, misc.CRAFTING_TABLE_RECIPE_SELECTED)
		if u == nil {
			break
		}
		var chrctr *entity.Entity
		var option string
		var cbData misc.CallbackData
		if u.Message != nil {
			chrctr = hs.characterHelper.GetCharacter(u.Message.From.ID)
			option = u.Message.Text
		} else {
			chrctr = hs.characterHelper.GetCharacter(u.CallbackQuery.From.ID)
			cbData = misc.CallbackData{}
			cbData.FromJSON(u.CallbackQuery.Data)
			option = cbData.Action
		}
		if chrctr == nil {
			continue
		}

		switch option {
		case bot.HIDEOUT_CRAFTING_TABLE, misc.CRAFTING_TABLE_RECIPE_SELECTED:
			player_C := chrctr.GetComponent("PlayerComponent").(*components.PlayerComponent)
			inventory_C := chrctr.GetComponent("InventoryComponent").(*components.InventoryComponent)
			recipes := inventory_C.GetItems(items.RECIPE)
			kb := misc.TGInlineKeyboard{Columns: 2}
			text := resources.HIDEOUT_CRAFTING_TABLE_TEXT
			if len(recipes) == 0 {
				kb.AddButton("You have no recipes", "")
			}
			if u.CallbackQuery != nil {
				itm_id, _ := strconv.Atoi(cbData.Payload)
				for _, r := range recipes {
					if r.ID == itm_id {
						item := r.GetItem().(*items.Recipe)
						text = fmt.Sprintf("%s\nRequires:\n", item.GetName())
						for _, itm := range item.Ingridients {
							text += fmt.Sprintf("%s (%d)\n", itm.GetItem().GetName(), itm.Qty)
						}
					}
				}
			}
			for _, r := range recipes {
				id_str := strconv.Itoa(r.ID)
				cbData := misc.CallbackData{Action: misc.CRAFTING_TABLE_RECIPE_SELECTED, Payload: id_str}
				kb.AddButton(fmt.Sprintf("%s (%d)", r.GetItem().GetName(), r.Qty), cbData.JSON())
			}

			if u.CallbackQuery != nil {
				msg := tgbotapi.NewEditMessageText(player_C.ChatID, u.CallbackQuery.Message.MessageID, text)
				msg.ReplyMarkup = kb.Generate()
				bot.GetInstance().GetBot().Send(msg)
			} else {
				msg := tgbotapi.NewMessage(player_C.ChatID, text)
				msg.ReplyMarkup = kb.Generate()
				bot.GetInstance().GetBot().Send(msg)
			}
		}
	}
}
