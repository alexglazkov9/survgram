package manager

import (
	"log"

	"github.com/alexglazkov9/survgram/activity"
	"github.com/alexglazkov9/survgram/components"
	"github.com/alexglazkov9/survgram/lootmanager"
	"github.com/alexglazkov9/survgram/misc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type ActivityManager struct {
	Activities []*activity.Activity

	LootManager *lootmanager.LootManager
}

func (a *ActivityManager) Update(dt float64) {
	for i, activity := range a.Activities {
		if activity.IsComplete {
			a.Activities = append(a.Activities[:i], a.Activities[i+1:]...)
			log.Printf("activities len: %d", len(a.Activities))
			continue
		}
		activity.Update(dt)
	}
}

func (a *ActivityManager) Add(ac *activity.Activity) {
	ac.LootManager = a.LootManager
	a.Activities = append(a.Activities, ac)
	log.Printf("activities len: %d", len(a.Activities))
}

func (a *ActivityManager) HandleInput(update tgbotapi.Update) {
	for _, act := range a.Activities {
		if act.Host.GetComponent("PlayerComponent").(*components.PlayerComponent).TelegramID == update.CallbackQuery.From.ID {
			cbData := misc.CallbackData{}
			cbData.FromJSON(update.CallbackQuery.Data)
			switch cbData.Action {
			case misc.ACTIVITY_CONTINUE:
				act.IsReadyForNext = true
			case misc.ACTIVITY_LEAVE:
				act.IsComplete = true
			}
			continue
		}
	}
}
