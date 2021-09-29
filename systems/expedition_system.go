package systems

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/alexglazkov9/survgram/activities"
	"github.com/alexglazkov9/survgram/bot"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/entity/enemies"
	"github.com/alexglazkov9/survgram/interfaces"
	"github.com/alexglazkov9/survgram/items"
	"github.com/alexglazkov9/survgram/misc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	ACTIVITIES_OPTIONS_NUMBER = 3
)

type ExpeditionSystem struct {
	currentDt float64

	manager         *entity.Manager
	characterHelper interfaces.CharacterHelper
}

func NewExpeditionSystem(manager *entity.Manager, characterHelper interfaces.CharacterHelper) *ExpeditionSystem {
	sys := &ExpeditionSystem{manager: manager, characterHelper: characterHelper}

	return sys
}

func (es *ExpeditionSystem) Update(dt float64) {
	/*Handle input*/
	for {
		u := bot.GetInstance().PopUpdate(
			misc.EXPEDITION_ACTIVITY_SELECTED,
			misc.EXPEDITION_CONTINUE,
			misc.EXPEDITION_LEAVE,
		)
		if u == nil {
			break
		}
		cbData := misc.CallbackData{}
		if u.CallbackQuery != nil {
			cbData.FromJSON(u.CallbackQuery.Data)
		}

		id, _ := strconv.Atoi(cbData.ID)
		ent := es.manager.GetEntityByID(id)
		expedition_C := ent.GetComponent("ExpeditionComponent").(*components.ExpeditionComponent)

		switch cbData.Action {
		case misc.EXPEDITION_ACTIVITY_SELECTED:
			if expedition_C.CurrentActivity != nil { //ignore input if user is already in the activity
				continue
			}

			selectedActivityId, _ := strconv.Atoi(cbData.Payload)
			act := expedition_C.GeneratedActivities[selectedActivityId]
			log.Println(act.Type)
			switch act.Type {
			case components.BATTLE_ACTIVITY:
				chrctr := es.characterHelper.GetCharacter(u.CallbackQuery.From.ID)
				player_C := chrctr.GetComponent("PlayerComponent").(*components.PlayerComponent)
				//TODO fix hardcoded 0 if makes sense
				spawnee := act.Spawnees[0]
				e := enemies.GetInstance().GetEnemyById(spawnee.Id)

				battleC := &components.SharedBattleComponent{
					CurrentState: components.PREACTIVITY,
					Messages:     make(map[int]tgbotapi.Message),
				}
				battleC.AddEnemy(&e)
				battleC.AddPlayer(chrctr)
				battleC.Messages[player_C.TelegramID] = expedition_C.Messages[player_C.TelegramID]

				statusC := &components.ActivityStatusComponent{
					IsComplete: false,
				}

				battle := es.manager.NewEntity()
				battle.AddComponent(battleC)
				battle.AddComponent(statusC)

				expedition_C.CurrentActivity = battle
				expedition_C.State = components.ACTIVITY_RUNNING
			case components.GATHERING_ACTIVITY:
				chrctr := es.characterHelper.GetCharacter(u.CallbackQuery.From.ID)
				spawnee := act.Spawnees[0]
				gatheringActivityC := &components.GatheringActivityComponent{
					CurrentState: components.PREACTIVITY,
					IsSuccesful:  true,
					Messages:     make(map[int]tgbotapi.Message),
					Players:      make([]*entity.Entity, 0),
					Resource: items.ItemBundle{
						ID:  spawnee.Id,
						Qty: spawnee.Qty,
					},
				}
				gatheringActivityC.Players = append(gatheringActivityC.Players, chrctr)
				statusC := &components.ActivityStatusComponent{
					IsComplete: false,
				}

				gatheringActivity := es.manager.NewEntity()
				gatheringActivity.AddComponent(gatheringActivityC)
				gatheringActivity.AddComponent(statusC)

				expedition_C.CurrentActivity = gatheringActivity
				expedition_C.State = components.ACTIVITY_RUNNING
			}
		case misc.EXPEDITION_CONTINUE:
			expedition_C.State = components.STARTING
		case misc.EXPEDITION_LEAVE:
			//Remove player and end expedition if no other players left
			for i, p := range expedition_C.Players {
				player_C := p.GetComponent("PlayerComponent").(*components.PlayerComponent)
				if player_C.TelegramID == u.CallbackQuery.From.ID {
					expedition_C.Players = append(expedition_C.Players[:i], expedition_C.Players[i+1:]...)
					if p.HasComponent("PlayerActivityComponent") {
						p.RemoveComponent("PlayerActivityComponent")
					}
					editTxt := tgbotapi.NewEditMessageText(player_C.ChatID, expedition_C.Messages[player_C.TelegramID].MessageID, "You have finished the expedition")
					bot.GetInstance().GetBot().Send(editTxt)
					//editKb := tgbotapi.NewEditMessageReplyMarkup(player_C.ChatID, expedition_C.Messages[player_C.TelegramID].MessageID, )
					if len(expedition_C.Players) == 0 {
						es.manager.DisposeEntity(ent)
					}
					SendMainMenuKeyboard(p)
				}
			}
		}
	}

	/*Handle logic*/
	for _, e := range es.manager.QueryEntities("ExpeditionComponent") {
		expedition_C := e.GetComponent("ExpeditionComponent").(*components.ExpeditionComponent)
		switch expedition_C.State {
		case components.STARTING:
			for _, p := range expedition_C.Players {
				player_C := p.GetComponent("PlayerComponent").(*components.PlayerComponent)
				//Send message with continue and leave btns
				loc := expedition_C.Location.(activities.Location)
				expedition_C.GeneratedActivities = nil
				tgkb := misc.TGInlineKeyboard{Columns: 1}
				for i := 0; i < ACTIVITIES_OPTIONS_NUMBER; i++ {
					pos_act := loc.PossibleActivities[rand.Intn(len(loc.PossibleActivities))]
					activity_config := components.ActivityConfig{
						Type:     pos_act.Type,
						Spawnees: []components.SpawneeConfig{GetSpawneeCfg(pos_act.SpawnChances)},
					}
					expedition_C.GeneratedActivities = append(expedition_C.GeneratedActivities, activity_config)
					cbData := misc.CallbackData{Action: misc.EXPEDITION_ACTIVITY_SELECTED, ID: strconv.Itoa(e.GetID()), Payload: fmt.Sprint(i)}
					tgkb.AddButton(string(activity_config.Type), cbData.JSON())
				}
				msg := tgbotapi.NewMessage(player_C.ChatID, "Select the activity: ")
				msg.ReplyMarkup = tgkb.Generate()
				expedition_C.Messages[player_C.TelegramID], _ = bot.GetInstance().GetBot().Send(msg)
				expedition_C.State = components.WAITING_INPUT
			}
		case components.ACTIVITY_RUNNING:
			status_C := expedition_C.CurrentActivity.GetComponent("ActivityStatusComponent").(*components.ActivityStatusComponent)
			if status_C.IsComplete {
				es.manager.DisposeEntity(expedition_C.CurrentActivity)
				expedition_C.CurrentActivity = nil
				expedition_C.State = components.WAITING_INPUT
				for _, p := range expedition_C.Players {
					player_C := p.GetComponent("PlayerComponent").(*components.PlayerComponent)

					/*Sends a message in-between activities*/
					msg := tgbotapi.NewMessage(player_C.ChatID, "Have a break between activities. You can access and manage your inventory, character and skills now!")
					tgkb := misc.TGInlineKeyboard{Columns: 2}
					e_id := strconv.Itoa(e.GetID())
					cbData := misc.CallbackData{Action: misc.EXPEDITION_CONTINUE, ID: e_id}
					tgkb.AddButton("Continue", cbData.JSON())
					cbData = misc.CallbackData{Action: misc.EXPEDITION_LEAVE, ID: e_id}
					tgkb.AddButton("Leave", cbData.JSON())
					msg.ReplyMarkup = tgkb.Generate()
					expedition_C.Messages[player_C.TelegramID], _ = bot.GetInstance().GetBot().Send(msg)
				}

			}
		}
	}
}

/*Helper function that returns id of the mob to spawn from the list of probabilities*/
func GetSpawneeCfg(spawn_chances []activities.SpawnChance) components.SpawneeConfig {
	cmltv := make([]float64, len(spawn_chances)) //cumulative
	for i, sc := range spawn_chances {
		if i == 0 {
			cmltv[i] = sc.Chance
			continue
		}
		cmltv[i] = cmltv[i-1] + sc.Chance
	}

	rnd_f := rand.Float64()
	for i, v := range cmltv {
		if rnd_f <= v {

			return components.SpawneeConfig{
				Id:  spawn_chances[i].Id,
				Qty: 1,
			}
		}
	}

	return components.SpawneeConfig{}
}
