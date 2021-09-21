package systems

import (
	"fmt"
	"log"

	"github.com/alexglazkov9/survgram/activities"
	"github.com/alexglazkov9/survgram/bot"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/interfaces"
	"github.com/alexglazkov9/survgram/misc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

/* Handles menu at the bottom of the screen built with ReplyKeyboard */
type MenuSystem struct {
	manager         *entity.Manager
	characterHelper interfaces.CharacterHelper
}

func NewMenuSystem(manager *entity.Manager, characterHelper interfaces.CharacterHelper) *MenuSystem {
	sys := &MenuSystem{manager: manager, characterHelper: characterHelper}

	return sys
}

func (ms *MenuSystem) Update(dt float64) {
	for {
		u := bot.GetInstance().PopUpdate(
			bot.MENU_OPTION_CHARACTER,
			bot.MENU_OPTION_EXPEDITION,
			bot.MENU_OPTION_NPCS,
			bot.MENU_OPTION_HIDEOUT,
			bot.MENU_OPTION_MAP,
			bot.MENU_BACK,

			bot.EXPEDITION_START,
		)
		if u == nil {
			break
		}
		chrctr := ms.characterHelper.GetCharacter(u.Message.From.ID)
		if chrctr == nil {
			//TODO Handle missing character
			continue
		}

		var menu_C *components.MenuComponent
		if chrctr.HasComponent("MenuComponent") {
			menu_C = chrctr.GetComponent("MenuComponent").(*components.MenuComponent)
		} else {
			menu_C = &components.MenuComponent{}
			menu_C.Menus.Push(GetMainMenu(chrctr))
			chrctr.AddComponent(menu_C)
		}

		var msg tgbotapi.Chattable
		log.Println(u.Message.Text)
		switch u.Message.Text {
		//MAP
		case bot.MENU_OPTION_MAP:
			player_C, _ := chrctr.GetComponent("PlayerComponent").(*components.PlayerComponent)

			//Add destinations to the keyboard
			loc := activities.GetLocations().GetLocation(player_C.CurrentLocation)
			kb := misc.TGInlineKeyboard{Columns: 2}
			for _, dest := range loc.Destinations {
				cbData := misc.CallbackData{Action: misc.GO_TO, Payload: fmt.Sprint(dest.GetID())}
				kb.AddButton(dest.Name, cbData.JSON())
			}
			msgT := tgbotapi.NewMessage(
				player_C.ChatID,
				"This is a map of Survgram. Nice image of the map with player's current position highlighted.\n\nClick a destination below to go there.",
			)
			msgT.ReplyMarkup = kb.Generate()
			msg = msgT
			bot.GetInstance().GetBot().Send(msg)
			continue //beacuse we don't want map message in the menu stack
		//START EXPEDITION
		case bot.EXPEDITION_START:
			player_C := chrctr.GetComponent("PlayerComponent").(*components.PlayerComponent)
			loc := activities.GetLocations().GetLocation(player_C.CurrentLocation)
			if chrctr.HasComponent("PlayerActivityComponent") { //check if player is busy and delete the message
				msg := tgbotapi.NewDeleteMessage(player_C.ChatID, u.Message.MessageID)
				bot.GetInstance().GetBot().DeleteMessage(msg)
				continue
			}

			//Start expedition
			expdtnComp := &components.ExpeditionComponent{
				State:    components.STARTING,
				Players:  []*entity.Entity{},
				IsReady:  true,
				Location: loc,
				Messages: make(map[int]tgbotapi.Message),
			}
			expdtn := ms.manager.NewEntity()
			expdtnComp.AddPlayer(chrctr)
			expdtn.AddComponent(expdtnComp)

			msg = GetExpeditionQuickMenu(chrctr)
		case bot.MENU_OPTION_EXPEDITION:
			msg = GetExpeditionMenu(chrctr)
		case bot.MENU_OPTION_HIDEOUT:
			msg = GetHideoutMenu(chrctr)
		case bot.MENU_BACK:
			var ok bool
			if len(menu_C.Menus) > 1 {
				menu_C.Menus.Pop()
			}
			msg, ok = menu_C.Menus.Top()
			if !ok {
				continue
			}
		}
		bot.GetInstance().GetBot().Send(msg)
		menu_C.Menus.Push(msg)

	}
}
