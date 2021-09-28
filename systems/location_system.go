package systems

import (
	"fmt"
	"strconv"

	"github.com/alexglazkov9/survgram/activities"
	"github.com/alexglazkov9/survgram/bot"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/interfaces"
	"github.com/alexglazkov9/survgram/misc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type LocationSystem struct {
	manager *entity.Manager

	characterHelper interfaces.CharacterHelper
}

func NewLocationSystem(manager *entity.Manager, characterHelper interfaces.CharacterHelper) *LocationSystem {
	sys := &LocationSystem{manager: manager, characterHelper: characterHelper}

	return sys
}

func (ls *LocationSystem) Update(dt float64) {
	for {
		u := bot.GetInstance().PopUpdate(misc.GO_TO)
		if u == nil {
			break
		}

		chrctr := ls.characterHelper.GetCharacter(u.CallbackQuery.From.ID)
		if chrctr == nil || chrctr.HasComponent("PlayerActivityComponent") { //check if player is busy
			//TODO Handle missing character
			continue
		}

		cbData := misc.CallbackData{}
		cbData.FromJSON(u.CallbackQuery.Data)
		loc_id, _ := strconv.Atoi(cbData.Payload)
		loc := activities.GetLocations().GetLocation(loc_id)
		location_C := chrctr.GetComponent("PlayerLocationComponent").(*components.PlayerLocationComponent)
		location_C.Destination = &loc_id
		location_C.TravelTime = 0
		chrctr.AddComponent(&components.PlayerActivityComponent{Type: "Travel"})
		deleteMsg := tgbotapi.NewDeleteMessage(u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.MessageID)
		bot.GetInstance().GetBot().Send(deleteMsg)
		msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, fmt.Sprintf("You are traveling to %s", loc.Name))
		msg.ReplyMarkup = GetQuickMenuKeyboard()
		bot.GetInstance().GetBot().Send(msg)
	}

	for _, e := range ls.manager.QueryEntities("PlayerLocationComponent") {
		location_C := e.GetComponent("PlayerLocationComponent").(*components.PlayerLocationComponent)
		if location_C.Destination == nil { //Player is not traveling
			continue
		}

		location_C.TravelTime += dt
		if location_C.TravelTime >= 5 {
			location_C.CurrentLocation = *location_C.Destination
			location_C.Destination = nil
			e.RemoveComponent("PlayerActivityComponent")
			ls.characterHelper.UpdateCharacter(e)
			menu_C := e.GetComponent("MenuComponent").(*components.MenuComponent)
			msg := GetMainMenu(e).(components.Menu)
			menu_C.Menus.Clear()
			menu_C.Menus.Push(msg)
			bot.GetInstance().GetBot().Send(msg.Msg)
		}
	}
}
