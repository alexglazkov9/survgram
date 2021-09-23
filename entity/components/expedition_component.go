package components

import (
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type ActivityType string

const (
	FISHING_ACTIVITY   ActivityType = "fishing"
	BATTLE_ACTIVITY    ActivityType = "battle"
	GATHERING_ACTIVITY ActivityType = "gathering"
)

type SpawneeConfig struct {
	Id  int
	Qty int
}

type ActivityConfig struct {
	Type     ActivityType
	Spawnees []SpawneeConfig
}

type ExpeditionState string

const (
	STARTING         ExpeditionState = "starting"
	WAITING_INPUT    ExpeditionState = "waiting_input"
	ACTIVITY_RUNNING ExpeditionState = "activity_running"
)

type ExpeditionComponent struct {
	BaseComponent `bson:"-" json:"-"`

	Players             []*entity.Entity
	Messages            map[int]tgbotapi.Message //Holds tgbotapi.Message sent to players in map[TelegramID_of_the_player]players_message in the current activity
	IsReady             bool
	GeneratedActivities []ActivityConfig
	State               ExpeditionState
	CurrentActivity     *entity.Entity

	Location interfaces.ILocation
}

func (ec *ExpeditionComponent) Clone() entity.IComponent {
	copy := *ec
	return &copy
}

func (ec *ExpeditionComponent) AddPlayer(p *entity.Entity) {
	p_activityC := &PlayerActivityComponent{
		Type: "expedition",
	}
	p.AddComponent(p_activityC)
	ec.Players = append(ec.Players, p)
}
