package components

import (
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/items"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type GatheringActivityComponent struct {
	BaseComponent `bson:"-" json:"-"`

	IsActivityComplete bool
	Players            []*entity.Entity
	Messages           map[int]tgbotapi.Message //Holds tgbotapi.Message sent to players in map[TelegramID_of_the_player]players_message in the current activity
	CurrentState       ActivityState
	Resource           items.ItemBundle
	Count              int
}
