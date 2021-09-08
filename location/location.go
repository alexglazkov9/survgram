package location

import (
	"github.com/alexglazkov9/survgram/activities"
	"github.com/alexglazkov9/survgram/battle"
	"github.com/alexglazkov9/survgram/lootmanager"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	ID                 primitive.ObjectID `bson:"_id"`
	Name               string
	Destinations       []primitive.ObjectID
	PossibleActivities []activities.ActivityConfig
}

func (l Location) GetActivity(bot *tgbotapi.BotAPI, lm *lootmanager.LootManager) activities.IActivity {
	var act activities.IActivity
	//switch l.PossibleActivities[rand.Intn(len(l.PossibleActivities))].ActivityType {
	switch activities.BATTLE {
	case activities.FISHING:
	case activities.BATTLE:
		act = battle.NewBattle(bot, lm)
	case activities.GATHERING:
	}
	return act
}

func (l Location) GetID() primitive.ObjectID {
	return l.ID
}

type LocationModel struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string
	Destinations []Location
}

// New - Creates a new character
func New(name string) *Location {
	return &Location{ID: primitive.NewObjectID(), Name: name}
}
