package expedition

import (
	"log"
	"strconv"

	charactermanager "github.com/alexglazkov9/survgram/character/manager"
	"github.com/alexglazkov9/survgram/components"
	"github.com/alexglazkov9/survgram/lootmanager"
	"github.com/alexglazkov9/survgram/misc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Expeditions struct {
	Expeditions []*Expedition

	LootManager      *lootmanager.LootManager
	CharacterManager *charactermanager.CharacterManager
}

func (a *Expeditions) Update(dt float64) {
	for i, journey := range a.Expeditions {
		if journey.IsComplete {
			a.CharacterManager.UpdateCharacter(journey.Host)
			a.Expeditions = append(a.Expeditions[:i], a.Expeditions[i+1:]...)
			log.Printf("activities len: %d", len(a.Expeditions))
			continue
		}
		journey.Update(dt)
	}
}

func (a *Expeditions) Add(ac *Expedition) {
	ac.LootManager = a.LootManager
	a.Expeditions = append(a.Expeditions, ac)
	ac.sendNextKeyboard()
	log.Printf("activities len: %d", len(a.Expeditions))
}

func (a *Expeditions) HandleInput(update tgbotapi.Update) {
	log.Println("Handling input")
	for _, act := range a.Expeditions {
		if act.Host.GetComponent("PlayerComponent").(*components.PlayerComponent).TelegramID == update.CallbackQuery.From.ID {
			cbData := misc.CallbackData{}
			cbData.FromJSON(update.CallbackQuery.Data)
			switch cbData.Action {
			case misc.EXPEDITION_CONTINUE:
				act.IsReadyForNext = true
			case misc.EXPEDITION_LEAVE:
				act.IsComplete = true
			case misc.ACTIVITY_SELECTED:
				i, _ := strconv.Atoi(cbData.Payload)
				act.SetSelectedActivity(i)
			}
			continue
		}
	}
}
