package activity

import (
	"log"

	"github.com/alexglazkov9/survgram/battle"
	"github.com/alexglazkov9/survgram/components"
	"github.com/alexglazkov9/survgram/entity"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Activity struct {
	Bot         tgbotapi.BotAPI
	host        *entity.Entity
	nextNow     bool
	nextNowTime float64
	curr        int
	queue       [4]byte

	battle battle.Battle

	IsComplete bool
}

func New(bot tgbotapi.BotAPI, char *entity.Entity) *Activity {
	return &Activity{Bot: bot, host: char, queue: [4]byte{'t', 'm', 't', 'm'}, curr: 0, nextNowTime: 0, nextNow: true}
}

func (a *Activity) Update(dt float64) {
	a.nextNowTime += dt

	if a.curr >= len(a.queue) && a.battle.Ended {
		a.IsComplete = true
	}

	if a.battle.Ended {
		log.Println("ended")
		a.battle = battle.Battle{}
		a.nextNow = true
	} else {
		a.battle.Update(dt)
	}

	if a.nextNow && a.nextNowTime >= 1 && a.curr < len(a.queue) {

		a.nextNow = false
		a.nextNowTime = 0
		switch a.queue[a.curr] {
		case 't':
			player_C := a.host.GetComponent("PlayerComponent").(*components.PlayerComponent)
			msg := tgbotapi.NewMessage(player_C.ChatID, "You are slowly approaching the Old Forrest!")
			a.Bot.Send(msg)
			a.nextNow = true
		case 'm':
			enem := entity.New()
			enem.AddComponent(&components.NameComponent{
				Parent: enem,
				Name:   "Sobaka",
			})
			enem.AddComponent(&components.AttackComponent{
				Parent:       enem,
				AttackDamage: 1,
			})
			enem.AddComponent(&components.HealthComponent{
				Parent:          enem,
				MaxHealthPoints: 150,
				HealthPoints:    150,
			})
			enem.AddComponent(&components.EffectsComponent{Parent: enem})

			a.battle = battle.Battle{Bot: a.Bot}
			a.battle.AddToEnemies(enem)
			a.battle.AddToPlayers(a.host)
			a.battle.Start()
		}
		a.curr++
	}
}
