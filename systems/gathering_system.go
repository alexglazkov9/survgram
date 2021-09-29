package systems

import (
	"strconv"

	"github.com/alexglazkov9/survgram/bot"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/interfaces"
	"github.com/alexglazkov9/survgram/misc"
)

/*
	System that controls all gathering activites in the game
*/
type GatheringSystem struct {
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
		gatheringActivity_C := gatheringActivity_e.GetComponent("GatheringActivityComponent").(*components.GatheringActivityComponent)

		switch cbData.Action {
		case misc.GATHERING_CORRECT:
			if gatheringActivity_C.CurrentState == components.RUNNING {
				gatheringActivity_C.Count++
				if gatheringActivity_C.Count >= 4 {
					// Gathered correctly - finish activity - update all messages
					gatheringActivity_C.CurrentState = components.POSTACTIVITY
					gatheringActivity_C.IsSuccesful = true
					for _, p := range gatheringActivity_C.Players {
						gatheringActivity_C.SendUpdate(p)
						//Distributre loot
						lootDispenser_C := &components.LootDispenserComponent{
							State: components.ADDED,
						}
						lootDispenser_C.AddItems(gatheringActivity_C.Resource)
						p.AddComponent(lootDispenser_C)
					}
				} else {
					// Gathered correctly - update message for the player
					chrctr := gs.characterHelper.GetCharacter(u.CallbackQuery.From.ID)
					gatheringActivity_C.SendUpdate(chrctr)
				}
			}
		case misc.GATHERING_INCORRECT:
			// Gathered incorrectly - finsih activity - update messages for all players
			gatheringActivity_C.CurrentState = components.POSTACTIVITY
			gatheringActivity_C.IsSuccesful = false
			for _, p := range gatheringActivity_C.Players {
				gatheringActivity_C.SendUpdate(p)
			}
			activity_status_c := gatheringActivity_e.GetComponent("ActivityStatusComponent").(*components.ActivityStatusComponent)
			activity_status_c.IsComplete = true
		}
	}

	/* Update all gathering activities */
	for _, entity := range gs.manager.QueryEntities("GatheringActivityComponent") {
		gatheringActivity_C := entity.GetComponent("GatheringActivityComponent").(*components.GatheringActivityComponent)

		switch gatheringActivity_C.CurrentState {
		case components.PREACTIVITY:
			for _, p := range gatheringActivity_C.Players {
				gatheringActivity_C.SendUpdate(p)
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
