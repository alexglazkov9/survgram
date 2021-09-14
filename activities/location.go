package activities

import (
	"math/rand"

	"github.com/alexglazkov9/survgram/items/loot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Location struct {
	ID                 int
	Name               string
	Destinations       []*Location
	PossibleActivities []ActivityConfig
}

/* Creates activity from one of the ActivityConfig assigned to this location and returns it */
func (l Location) GetActivity(bot *tgbotapi.BotAPI, lm *loot.LootDispenser) IActivity {
	var act IActivity
	act_index := rand.Intn(len(l.PossibleActivities))
	switch l.PossibleActivities[act_index].Type {
	case FISHING:
	case BATTLE:
		act = NewBattle(bot, lm, l.PossibleActivities[act_index].SpawnChances)
	case GATHERING:
		act = NewGatheringActivity(bot, lm, l.PossibleActivities[act_index].SpawnChances)
	}
	return act
}

func (l Location) GetID() int {
	return l.ID
}

// NewLocation - Creates a new character
func NewLocation(id int, name string) *Location {
	return &Location{ID: id, Name: name, Destinations: make([]*Location, 0)}
}
