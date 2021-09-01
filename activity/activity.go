package activity

import (
	"log"

	"github.com/alexglazkov9/survgram/battle"
	"github.com/alexglazkov9/survgram/character"
	"github.com/alexglazkov9/survgram/components"
	"github.com/alexglazkov9/survgram/enemy"
	"github.com/alexglazkov9/survgram/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Activity struct {
	Bot         tgbotapi.BotAPI
	host        *character.Character
	nextNow     bool
	nextNowTime float64
	curr        int
	queue       [4]byte

	battle battle.Battle

	IsComplete bool
}

func New(bot tgbotapi.BotAPI, char *character.Character) *Activity {
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

	if a.nextNow && a.nextNowTime >= 3 && a.curr < len(a.queue) {

		a.nextNow = false
		switch a.queue[a.curr] {
		case 't':
			msg := tgbotapi.NewMessage(a.host.ChatID, "You are slowly approaching the Old Forrest!")
			a.Bot.Send(msg)
			a.nextNow = true
		case 'm':
			var enem interfaces.Battler
			enem = &enemy.Enemy{
				NameComponent: components.NameComponent{
					Name: "Sobaka",
				},
				AttackComponent: components.AttackComponent{
					AttackDamage: 1,
				},
				HealthComponent: components.HealthComponent{
					MaxHealthPoints: 10,
					HealthPoints:    10,
				},
			}
			var player interfaces.Battler = a.host
			a.battle = battle.Battle{Bot: a.Bot}
			a.battle.AddToEnemies(enem)
			a.battle.AddToPlayers(player)
			a.battle.Start()
		}
		a.curr++
	}
}

// func (a Activity) Start() {
// 	curr := 0
// 	queue := []byte{'t', 'm', 't'}
// 	for curr <= len(queue) {
// 		switch queue[curr] {
// 		case 't':
// 			msg := tgbotapi.NewMessage(a.unit.ChatID, "You are slowly approaching the Old Forrest!")
// 			a.Bot.Send(msg)
// 		case 'm':
// 			enem := enemy.Enemy{Name: "Sobaka", Health: 10, Attack: 1}
// 			a.StartBattle(enem)
// 		}
// 		curr++
// 	}
// }

func (a Activity) StartBattle(enemy enemy.Enemy) {

}
