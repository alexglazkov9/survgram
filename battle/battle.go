package battle

import (
	"fmt"
	"log"

	"github.com/alexglazkov9/survgram/components"
	"github.com/alexglazkov9/survgram/entity"

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
	enemies []*entity.Entity
	players []*entity.Entity

	battleLog []string

	nextUpdateTimer float64
	messages        map[int]tgbotapi.Message

	rewards map[int]string

	Ended bool
}

func (b *Battle) Update(dt float64) {
	//Do i need this?
	if b.Ended {
		return
	}

	for _, e := range b.getAllEntities() {
		e.Update(dt)
	}

	for _, e := range b.GetAllEnemies(true) {

		ability_c, ok := e.GetComponent("AbilityComponent").(*components.AbilityComponent)
		if ok {
			if ability_c.IsReadyToCast() {
				ability_c.CastAbility(b.enemies, b.players)
			}
		}

		attack_c := e.GetComponent("AttackComponent").(*components.AttackComponent)
		if attack_c.IsReadyToAttack() {
			attackDetails := attack_c.Attack(b.players)
			attackDetails.Attacker = e.GetComponent("NameComponent").(*components.NameComponent).GetName()
			if ability_c != nil {
				ability_c.GainEnergy(attackDetails.Damage)
			}
			b.battleLog = append(b.battleLog, fmt.Sprintf("%s attacks %s for %d damage", attackDetails.Attacker, attackDetails.Target, attackDetails.Damage))
			if !b.isAnyoneAlive(PLAYERS) {
				log.Println("players lost")
				b.Ended = true
				b.sendUpdate()
			}
		}
	}

	for _, p := range b.GetAllPlayers(true) {

		ability_c, ok := p.GetComponent("AbilityComponent").(*components.AbilityComponent)
		if ok {
			if ability_c.IsReadyToCast() {
				ability_c.CastAbility(b.players, b.enemies)
			}
		}

		attack_c := p.GetComponent("AttackComponent").(*components.AttackComponent)
		if attack_c.IsReadyToAttack() {
			attackDetails := attack_c.Attack(b.enemies)
			attackDetails.Attacker = p.GetComponent("NameComponent").(*components.NameComponent).GetName()
			if ability_c != nil {
				ability_c.GainEnergy(attackDetails.Damage)
			}
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

		msg := tgbotapi.NewMessage(p.GetComponent("PlayerComponent").(*components.PlayerComponent).ChatID, b.generateStatusText())
		b.messages[p.GetComponent("PlayerComponent").(*components.PlayerComponent).TelegramID], _ = b.Bot.Send(msg)
	}
	b.Ended = false
}

func (b Battle) sendUpdate() {
	for _, p := range b.players {
		msg := tgbotapi.NewEditMessageText(
			p.GetComponent("PlayerComponent").(*components.PlayerComponent).ChatID,
			b.messages[p.GetComponent("PlayerComponent").(*components.PlayerComponent).TelegramID].MessageID,
			b.generateStatusText(),
		)
		b.Bot.Send(msg)
	}
}

func (b Battle) generateStatusText() string {
	var playerStatus string
	for _, p := range b.players {
		name_C := p.GetComponent("NameComponent").(*components.NameComponent)
		health_c := p.GetComponent("HealthComponent").(*components.HealthComponent)
		playerStatus = fmt.Sprintf("%s\nHealth: %d/%d\n", name_C.GetName(), health_c.GetHealth(), health_c.GetMaxHealth())
	}
	status := playerStatus
	for _, e := range b.enemies {
		name_C := e.GetComponent("NameComponent").(*components.NameComponent)
		health_c := e.GetComponent("HealthComponent").(*components.HealthComponent)
		playerStatus = fmt.Sprintf("\n%s\nHealth: %d/%d\n", name_C.GetName(), health_c.GetHealth(), health_c.GetMaxHealth())
	}
	status += (playerStatus + "\n")
	for _, s := range b.battleLog {
		status += (s + "\n")
	}
	return status
}

func (b Battle) getAllEntities() []*entity.Entity {
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

func (b Battle) GetAllPlayers(args ...bool) []*entity.Entity {
	if len(args) > 0 && args[0] {
		result := make([]*entity.Entity, 0)
		for _, p := range b.players {
			health_C := p.GetComponent("HealthComponent").(*components.HealthComponent)
			if health_C.IsAlive() {
				result = append(result, p)
			}
		}
		return result
	}

	return b.players
}

func (b Battle) GetAllEnemies(args ...bool) []*entity.Entity {
	if len(args) > 0 && args[0] {
		result := make([]*entity.Entity, 0)
		for _, e := range b.enemies {
			health_C := e.GetComponent("HealthComponent").(*components.HealthComponent)
			if health_C.IsAlive() {
				result = append(result, e)
			}
		}
		return result
	}

	return b.players
}

func (b *Battle) AddToEnemies(e *entity.Entity) {
	e.AddComponent(&components.LogComponent{
		Parent: e,
		Log:    &b.battleLog,
	})
	b.enemies = append(b.enemies, e)
}

func (b *Battle) AddToPlayers(e *entity.Entity) {
	e.AddComponent(&components.LogComponent{
		Parent: e,
		Log:    &b.battleLog,
	})
	b.players = append(b.players, e)
}
