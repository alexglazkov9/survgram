package activities

import (
	"github.com/alexglazkov9/survgram/entity"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type IActivity interface {
	Start()
	Update(dt float64)
	IsComplete() bool
	GetDescription() string
	AddPlayer(player *entity.Entity, msg ...*tgbotapi.Message)
}

type ActivityState string

const (
	PREACTIVITY  ActivityState = "preactivity"
	RUNNING      ActivityState = "running"
	POSTACTIVITY ActivityState = "postactivity"
)
