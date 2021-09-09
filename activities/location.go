package activities

import (
	"github.com/alexglazkov9/survgram/items/loot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Location struct {
	ID                 int
	Name               string
	Destinations       []*Location
	PossibleActivities []ActivityConfig
}

func (l Location) GetActivity(bot *tgbotapi.BotAPI, lm *loot.LootManager) IActivity {
	var act IActivity
	//switch l.PossibleActivities[rand.Intn(len(l.PossibleActivities))].ActivityType {
	switch BATTLE {
	case FISHING:
	case BATTLE:
		act = NewBattle(bot, lm)
	case GATHERING:
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
