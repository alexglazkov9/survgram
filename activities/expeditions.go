package activities

import (
	"log"
	"strconv"

	charactermanager "github.com/alexglazkov9/survgram/database"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/items/loot"
	"github.com/alexglazkov9/survgram/misc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Expeditions struct {
	Expeditions []*Expedition

	LootManager      *loot.LootDispenser
	CharacterManager *charactermanager.CharacterManager
}

/* Updates all existing expeditions and monitors
when expedition is complete to remove it */
func (a *Expeditions) Update(dt float64) {
	for i, expedition := range a.Expeditions {
		if expedition.IsComplete {
			a.CharacterManager.UpdateCharacter(expedition.Host)
			a.Expeditions = append(a.Expeditions[:i], a.Expeditions[i+1:]...)
			log.Printf("activities len: %d", len(a.Expeditions))
			continue
		}
		expedition.Update(dt)
	}
}

/* Adds expedition to the list of expeditions and starts it */
func (a *Expeditions) Add(ac *Expedition) {
	ac.LootManager = a.LootManager
	a.Expeditions = append(a.Expeditions, ac)
	ac.sendNextKeyboard()
	log.Printf("activities len: %d", len(a.Expeditions))
}

func (a *Expeditions) HandleInput(update tgbotapi.Update) {
	for _, expdtn := range a.Expeditions {
		if expdtn.Host.GetComponent("PlayerComponent").(*components.PlayerComponent).TelegramID == update.CallbackQuery.From.ID {
			cbData := misc.CallbackData{}
			cbData.FromJSON(update.CallbackQuery.Data)
			switch cbData.Action {
			case misc.EXPEDITION_CONTINUE:
				expdtn.IsReadyForNext = true
			case misc.EXPEDITION_LEAVE:
				expdtn.EndExpedition()
			case misc.ACTIVITY_SELECTED:
				i, _ := strconv.Atoi(cbData.Payload)
				expdtn.SetSelectedActivity(i)
			case misc.GATHERING_CORRECT, misc.GATHERING_INCORRECT:
				expdtn.HandleInput(update)
			}
			continue
		}
	}
}
