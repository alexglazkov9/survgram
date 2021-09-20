package systems

import (
	"fmt"
	"strings"

	"github.com/alexglazkov9/survgram/activities"
	"github.com/alexglazkov9/survgram/bot"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/interfaces"
	"github.com/alexglazkov9/survgram/resources"
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
		u := bot.GetInstance().PopUpdate(bot.MENU_OPTION_EXPEDITION, bot.MENU_OPTION_CHARACTER, bot.MENU_OPTION_NPCS, bot.EXPEDITION_START)
		if u == nil {
			break
		}

		chrctr := ls.characterHelper.GetCharacter(u.Message.From.ID)
		if chrctr == nil {
			//TODO Handle missing character
			continue
		}
		player_C := chrctr.GetComponent("PlayerComponent").(*components.PlayerComponent)
		loc := activities.GetLocations().GetLocation(player_C.CurrentLocation)
		switch u.Message.Text {
		case bot.MENU_OPTION_EXPEDITION:
			msgTxt := "Possible activities in this location:\n"
			for _, act := range loc.PossibleActivities {
				msgTxt += fmt.Sprintf("%.0f%% %s \n", act.ActivityChance*100, strings.Title(string(act.Type)))
			}
			kb := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(bot.EXPEDITION_START),
				tgbotapi.NewKeyboardButton(bot.EXPEDITION_ACTIVITIES_INFO),
				tgbotapi.NewKeyboardButton(bot.EXPEDITION_BACK_TO_MAIN_MENU),
			))
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, msgTxt)
			msg.ReplyMarkup = kb
			bot.GetInstance().GetBot().Send(msg)
		case bot.MENU_OPTION_NPCS:
		case bot.MENU_OPTION_CHARACTER:
		case bot.MENU_OPTION_HIDEOUT:

		case bot.EXPEDITION_START:
			if chrctr.HasComponent("PlayerActivityComponent") { //check if player is busy
				msg := tgbotapi.NewDeleteMessage(u.Message.Chat.ID, u.Message.MessageID)
				bot.GetInstance().GetBot().DeleteMessage(msg)
				continue
			}

			expdtnComp := &components.ExpeditionComponent{
				State:    components.STARTING,
				Players:  []*entity.Entity{},
				IsReady:  true,
				Location: loc,
				Messages: make(map[int]tgbotapi.Message),
			}
			expdtn := ls.manager.NewEntity()
			expdtnComp.AddPlayer(chrctr)
			expdtn.AddComponent(expdtnComp)
			msg := tgbotapi.NewMessage(player_C.ChatID, fmt.Sprintf("%s at %s", resources.EXPEDITION_ENTERED_TEXT, loc.Name))
			kb := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(bot.EXPEDITION_RUNNING_MENU_CHARACTER),
					tgbotapi.NewKeyboardButton(bot.EXPEDITION_RUNNING_MENU_INVENTORY),
					tgbotapi.NewKeyboardButton(bot.EXPEDITION_RUNNING_MENU_ABILITIES),
				),
			)
			msg.ReplyMarkup = kb
			bot.GetInstance().GetBot().Send(msg)
		}
	}
}
