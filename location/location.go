package location

import (
	"github.com/alexglazkov9/survgram/activity"
	"github.com/alexglazkov9/survgram/battle"
	"github.com/alexglazkov9/survgram/lootmanager"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	ID                 primitive.ObjectID `bson:"_id"`
	Name               string
	Destinations       []primitive.ObjectID
	PossibleActivities []activity.ActivityConfig
}

func (l Location) GetActivity(bot *tgbotapi.BotAPI, lm *lootmanager.LootManager) activity.IActivity {
	var act activity.IActivity
	//switch l.PossibleActivities[rand.Intn(len(l.PossibleActivities))].ActivityType {
	switch activity.BATTLE {
	case activity.FISHING:
	case activity.BATTLE:
		act = battle.NewBattle(bot, lm)
	case activity.GATHERING:
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
