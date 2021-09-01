package battle

import (
	"fmt"
	"log"

	"github.com/alexglazkov9/survgram/character"
	"github.com/alexglazkov9/survgram/interfaces"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Side int

const (
	PLAYERS Side = iota
	ENEMIES
)

const UPDATE_PERIOD float64 = 2

type Battle struct {
	Bot     tgbotapi.BotAPI
	enemies []interfaces.Battler
	players []interfaces.Battler

	battleLog []string

	nextUpdateTimer float64
	messages        map[int]tgbotapi.Message

	Ended bool
}

func (b *Battle) Update(dt float64) {
	if b.Ended {
		return
	}
	for _, e := range b.getAllEntities() {
		e.Update(dt)
	}
	for _, e := range b.enemies {
		if e.IsReadyToCast() {
			e.CastAbility()
		}
		if e.IsReadyToAttack() {
			attackDetails := e.Attack(b.players)
			attackDetails.Attacker = e.GetName()
			e.GainEnergy(attackDetails.Damage)
			b.battleLog = append(b.battleLog, fmt.Sprintf("%s attacks %s for %d damage", attackDetails.Attacker, attackDetails.Target, attackDetails.Damage))
			if !b.isAnyoneAlive(PLAYERS) {
				log.Println("players lost")
				b.Ended = true
				b.sendUpdate()
			}
		}
	}
	for _, e := range b.players {
		if e.IsReadyToAttack() {
			attackDetails := e.Attack(b.enemies)
			attackDetails.Attacker = e.GetName()
			b.battleLog = append(b.battleLog, fmt.Sprintf("%s attacks %s for %d damage", attackDetails.Attacker, attackDetails.Target, attackDetails.Damage))
			if !b.isAnyoneAlive(ENEMIES) {
				log.Println("players won")
				b.Ended = true
				b.sendUpdate()
			}
		}
	}

	b.nextUpdateTimer += dt
	if b.nextUpdateTimer >= UPDATE_PERIOD {
		b.nextUpdateTimer = 0
		b.sendUpdate()
	}
}

func (b *Battle) Start() {
	b.messages = make(map[int]tgbotapi.Message)
	b.battleLog = make([]string, 0)
	for _, p := range b.players {
		msg := tgbotapi.NewMessage(p.(*character.Character).ChatID, b.generateStatusText())
		b.messages[p.(*character.Character).TelegramID], _ = b.Bot.Send(msg)
	}
	b.Ended = false
}

func (b Battle) sendUpdate() {
	for _, p := range b.players {
		msg := tgbotapi.NewEditMessageText(p.(*character.Character).ChatID, b.messages[p.(*character.Character).TelegramID].MessageID, b.generateStatusText())
		b.Bot.Send(msg)
	}
}

func (b Battle) generateStatusText() string {
	var playerStatus string
	for _, p := range b.players {
		playerStatus = fmt.Sprintf("%s\nHealth: %d/%d\n", p.GetName(), p.GetHealth(), p.GetMaxHealth())
	}
	status := playerStatus
	for _, e := range b.enemies {
		playerStatus = fmt.Sprintf("\n%s\nHealth: %d/%d\n", e.GetName(), e.GetHealth(), e.GetMaxHealth())
	}
	status += (playerStatus + "\n")
	for _, s := range b.battleLog {
		status += (s + "\n")
	}
	return status
}

func (b Battle) getAllEntities() []interfaces.Battler {
	return append(b.enemies, b.players...)
}

func (b Battle) isAnyoneAlive(side Side) bool {
	if side == PLAYERS {
		return len(b.GetAllPlayers(true)) > 0
	}
	if side == ENEMIES {
		return len(b.GetAllEnemies(true)) > 0
	}
	return false
}

func (b Battle) GetAllPlayers(args ...bool) []interfaces.Battler {
	if len(args) > 0 && args[0] {
		result := make([]interfaces.Battler, 0)
		for _, p := range b.players {
			if p.IsAlive() {
				result = append(result, p)
			}
		}
		return result
	}

	return b.players
}

func (b Battle) GetAllEnemies(args ...bool) []interfaces.Battler {
	if len(args) > 0 && args[0] {
		result := make([]interfaces.Battler, 0)
		for _, e := range b.enemies {
			if e.IsAlive() {
				result = append(result, e)
			}
		}
		return result
	}

	return b.players
}

func (b *Battle) AddToEnemies(e interfaces.Battler) {
	b.enemies = append(b.enemies, e)
}

func (b *Battle) AddToPlayers(e interfaces.Battler) {
	b.players = append(b.players, e)
}
