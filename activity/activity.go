package activity

import (
	"github.com/alexglazkov9/survgram/battle"
	"github.com/alexglazkov9/survgram/character"
	"github.com/alexglazkov9/survgram/enemy"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Activity struct {
	Bot         tgbotapi.BotAPI
	unit        character.Character
	nextNow     bool
	nextNowTime float64
	curr        int
	queue       [3]byte

	battle battle.Battle
}

func New(bot tgbotapi.BotAPI, char character.Character) *Activity {
	return &Activity{Bot: bot, unit: char, queue: [3]byte{'t', 'm', 't'}, curr: 0, nextNowTime: 0, nextNow: true}
}

func (a *Activity) Update(dt float64) {
	a.nextNowTime += dt

	a.battle.Update(dt)

	if a.nextNow && a.nextNowTime >= 2 {

		a.nextNow = false
		switch a.queue[a.curr] {
		case 't':
			msg := tgbotapi.NewMessage(a.unit.ChatID, "You are slowly approaching the Old Forrest!")
			a.Bot.Send(msg)
			a.nextNow = true
		case 'm':
			enem := enemy.Enemy{AttackDamage: 10}
			a.battle = battle.Battle{}
			a.battle.AddToSideA(&enem)
			a.battle.AddToSideB(&a.unit)
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
