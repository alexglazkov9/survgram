package systems

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/alexglazkov9/survgram/bot"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/interfaces"
	"github.com/alexglazkov9/survgram/misc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type GatheringSystem struct {
	currentDt float64

	manager         *entity.Manager
	characterHelper interfaces.CharacterHelper
}

func NewGatheringSystem(manager *entity.Manager, characterHelper interfaces.CharacterHelper) *GatheringSystem {
	sys := &GatheringSystem{manager: manager, characterHelper: characterHelper}

	return sys
}

func (gs *GatheringSystem) Update(dt float64) {
	//Handle input
	for {
		u := bot.GetInstance().PopUpdate(
			misc.GATHERING_CORRECT,
			misc.GATHERING_INCORRECT,
		)
		if u == nil {
			break
		}
		cbData := misc.CallbackData{}
		cbData.FromJSON(u.CallbackQuery.Data)
		entity_id, _ := strconv.Atoi(cbData.ID)
		gatheringActivity_e := gs.manager.GetEntityByID(entity_id)
		log.Println(cbData.ID)
		log.Println(entity_id)
		log.Println(gatheringActivity_e)
		gatheringActivity_C := gatheringActivity_e.GetComponent("GatheringActivityComponent").(*components.GatheringActivityComponent)

		switch cbData.Action {
		case misc.GATHERING_CORRECT:
			if gatheringActivity_C.CurrentState == components.RUNNING {
				gatheringActivity_C.Count++
				if gatheringActivity_C.Count >= 4 {
					gatheringActivity_C.CurrentState = components.POSTACTIVITY
					for _, p := range gatheringActivity_C.Players {
						player_C := p.GetComponent("PlayerComponent").(*components.PlayerComponent)
						msg := gatheringActivity_C.Messages[player_C.TelegramID]
						edit_msg := tgbotapi.NewEditMessageText(
							msg.Chat.ID,
							msg.MessageID,
							fmt.Sprintf("You succesfully gathered %s", gatheringActivity_C.Resource.GetItem().GetName()),
						)
						bot.GetInstance().GetBot().Send(edit_msg)
						//Distributre loot
						lootDispenser_C := &components.LootDispenserComponent{
							State: components.ADDED,
						}
						lootDispenser_C.AddItems(gatheringActivity_C.Resource)
						p.AddComponent(lootDispenser_C)
					}

					//end activity
					//distribute loot
				} else {
					msgT := tgbotapi.NewEditMessageText(
						u.CallbackQuery.Message.Chat.ID,
						u.CallbackQuery.Message.MessageID,
						fmt.Sprintf("Pick the correct button to gather %d/%d", gatheringActivity_C.Count, 4),
					)
					kb := misc.TGInlineKeyboard{Columns: 2}
					btn_index := rand.Intn(4)
					for i := 0; i < 4; i++ {
						if btn_index == i {
							kb.AddButton(
								"Gather",
								misc.CallbackData{
									Action: misc.GATHERING_CORRECT,
									ID:     fmt.Sprint(gatheringActivity_C.Parent.ID),
								}.JSON(),
							)
						} else {
							kb.AddButton("-", misc.CallbackData{Action: misc.GATHERING_INCORRECT, ID: fmt.Sprint(gatheringActivity_C.Parent.ID)}.JSON())
						}
					}
					reply_markup := kb.Generate()
					msgT.ReplyMarkup = reply_markup
					bot.GetInstance().GetBot().Send(msgT)
				}
			}
		case misc.GATHERING_INCORRECT:
			//end activity
			gatheringActivity_C.CurrentState = components.POSTACTIVITY
			for _, p := range gatheringActivity_C.Players {
				player_C := p.GetComponent("PlayerComponent").(*components.PlayerComponent)
				name_C := p.GetComponent("NameComponent").(*components.NameComponent)
				msg := gatheringActivity_C.Messages[player_C.TelegramID]
				edit_msg := tgbotapi.NewEditMessageText(msg.Chat.ID, msg.MessageID, fmt.Sprintf("Gathering failed by %s", name_C.GetName()))
				bot.GetInstance().GetBot().Send(edit_msg)
			}
			activity_status_c := gatheringActivity_e.GetComponent("ActivityStatusComponent").(*components.ActivityStatusComponent)
			activity_status_c.IsComplete = true
		}
	}
	for _, entity := range gs.manager.QueryEntities("GatheringActivityComponent") {
		gatheringActivity_C := entity.GetComponent("GatheringActivityComponent").(*components.GatheringActivityComponent)

		switch gatheringActivity_C.CurrentState {
		case components.PREACTIVITY:
			for _, p := range gatheringActivity_C.Players {
				player_C := p.GetComponent("PlayerComponent").(*components.PlayerComponent)
				var msg tgbotapi.Chattable
				if _, ok := gatheringActivity_C.Messages[player_C.TelegramID]; ok {
					msgT := tgbotapi.NewEditMessageText(player_C.ChatID, gatheringActivity_C.Messages[player_C.TelegramID].MessageID, fmt.Sprintf("Pick the correct button to gather %d/%d", gatheringActivity_C.Count, 4))
					kb := misc.TGInlineKeyboard{Columns: 2}
					btn_index := rand.Intn(4)
					for i := 0; i < 4; i++ {
						if btn_index == i {
							kb.AddButton("Gather", misc.CallbackData{Action: misc.GATHERING_CORRECT, ID: fmt.Sprint(gatheringActivity_C.Parent.ID)}.JSON())
						} else {
							kb.AddButton("-", misc.CallbackData{Action: misc.GATHERING_INCORRECT, ID: fmt.Sprint(gatheringActivity_C.Parent.ID)}.JSON())
						}
					}
					reply_markup := kb.Generate()
					msgT.ReplyMarkup = reply_markup
					msg = msgT
				} else {
					log.Println("NEw")
					msgT := tgbotapi.NewMessage(player_C.ChatID, fmt.Sprintf("Pick the correct button to gather %d/%d", gatheringActivity_C.Count, 4))
					kb := misc.TGInlineKeyboard{Columns: 2}
					btn_index := rand.Intn(4)
					for i := 0; i < 4; i++ {
						if btn_index == i {
							kb.AddButton("Gather", misc.CallbackData{Action: misc.GATHERING_CORRECT, ID: fmt.Sprint(gatheringActivity_C.Parent.ID)}.JSON())
						} else {
							kb.AddButton("-", misc.CallbackData{Action: misc.GATHERING_INCORRECT, ID: fmt.Sprint(gatheringActivity_C.Parent.ID)}.JSON())
						}
					}
					reply_markup := *kb.Generate()
					msgT.ReplyMarkup = reply_markup
					msg = msgT
				}
				gatheringActivity_C.Messages[player_C.TelegramID], _ = bot.GetInstance().GetBot().Send(msg)
			}
			gatheringActivity_C.CurrentState = components.RUNNING
		case components.RUNNING:
		case components.POSTACTIVITY:
			// Check if anyone still looting and finish battle
			stillLooting := false
			for _, p := range gatheringActivity_C.Players {
				if p.HasComponent("LootDispenserComponent") {
					stillLooting = true
				}
			}

			if !stillLooting {
				entity.GetComponent("ActivityStatusComponent").(*components.ActivityStatusComponent).IsComplete = true
				for _, p := range gatheringActivity_C.Players {
					gs.characterHelper.UpdateCharacter(p)
				}
			}
		}
	}
}
