package components

import (
	"fmt"
	"math/rand"

	"github.com/alexglazkov9/survgram/bot"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/items"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type ActivityState string

const (
	PREACTIVITY  ActivityState = "preactivity"
	RUNNING      ActivityState = "running"
	POSTACTIVITY ActivityState = "postactivity"
)

type SharedBattleComponent struct {
	BaseComponent `bson:"-" json:"-"`

	IsBattleComplete bool
	Enemies          []*entity.Entity
	Players          []*entity.Entity
	Messages         map[int]tgbotapi.Message //Holds tgbotapi.Message sent to players in map[TelegramID_of_the_player]players_message in the current activity
	CurrentState     ActivityState
	BattleLog        []string
	NextUpdateTimer  float64
}

/* Returns list of all enemies(entities on the other side) of the entity holding the component in the battle,
pass true to return ALIVE enemies only */
func (bc SharedBattleComponent) GetEnemiesForEntity(e *entity.Entity, alive ...bool) []*entity.Entity {
	var result []*entity.Entity
	found := false
	for _, p := range bc.Players {
		if p.GetID() == e.GetID() {
			result = append(result, bc.Enemies...)
			found = true
			break
		}
	}

	if !found {
		result = append(result, bc.Players...)
	}

	if len(alive) > 0 && alive[0] {
		for i, e := range result {
			health_C := e.GetComponent("HealthComponent").(*HealthComponent)
			if !health_C.IsAlive() {
				result = append(result[:i], result[i+1:]...)
			}
		}
	}

	return result
}

/* Returns list of all allies(entities on the same side) of the entity holding the component in the battle,
pass true to return ALIVE allies only */
func (bc SharedBattleComponent) GetAlliesForEntity(e *entity.Entity, alive ...bool) []*entity.Entity {
	var result []*entity.Entity
	found := false
	for _, p := range bc.Players {
		if p.GetID() == e.GetID() {
			result = append(result, bc.Players...)
			found = true
			break
		}
	}

	if !found {
		result = append(result, bc.Enemies...)
	}

	if len(alive) > 0 && alive[0] {
		for i, e := range result {
			health_C := e.GetComponent("HealthComponent").(*HealthComponent)
			if !health_C.IsAlive() {
				result = append(result[:i], result[i+1:]...)
			}
		}
	}

	return result
}

/* Returns list of all allies of the entity holding the component in the battle, pass true to return ALIVE allies only */
func (bc SharedBattleComponent) GetAllEntities(alive ...bool) []*entity.Entity {
	var result []*entity.Entity
	if len(alive) > 0 && alive[0] {
		result = append(bc.Enemies, bc.Players...)
		for i, e := range result {
			health_C := e.GetComponent("HealthComponent").(*HealthComponent)
			if !health_C.IsAlive() {
				result = append(result[:i], result[i+1:]...)
			}
		}
	} else {
		result = append(bc.Enemies, bc.Players...)
	}

	return result
}

func (bc *SharedBattleComponent) HasAnyoneWon() bool {
	foundAlive := false
	for _, p := range bc.Players {
		health_C := p.GetComponent("HealthComponent").(*HealthComponent)
		if health_C.IsAlive() {
			foundAlive = true
			break
		}
	}

	if foundAlive {
		foundAlive = false
		for _, p := range bc.Enemies {
			health_C := p.GetComponent("HealthComponent").(*HealthComponent)
			if health_C.IsAlive() {
				return false
			}
		}
	}
	return true
}

func (bc *SharedBattleComponent) AddLogEntry(entry string) {
	bc.BattleLog = append(bc.BattleLog, entry)
}

func (bc SharedBattleComponent) SendUpdate() {
	for _, p := range bc.Players {
		player_C, player_C_ok := p.GetComponent("PlayerComponent").(*PlayerComponent)
		if player_C_ok {
			msg := tgbotapi.NewEditMessageText(
				player_C.ChatID,
				bc.Messages[player_C.TelegramID].MessageID,
				bc.generateStatusText(),
			)
			bot.GetInstance().GetBot().Send(msg)
		}
	}
}

func (bc *SharedBattleComponent) StartBattle() {
	bc.BattleLog = make([]string, 0)
	for _, p := range bc.Players {
		player_C := p.GetComponent("PlayerComponent").(*PlayerComponent)
		var msg tgbotapi.Chattable
		if _, ok := bc.Messages[player_C.TelegramID]; ok {
			msg = tgbotapi.NewEditMessageText(player_C.ChatID, bc.Messages[player_C.TelegramID].MessageID, bc.generateStatusText())
		} else {
			msg = tgbotapi.NewMessage(player_C.ChatID, bc.generateStatusText())
		}
		bc.Messages[player_C.TelegramID], _ = bot.GetInstance().GetBot().Send(msg)
	}
	bc.CurrentState = RUNNING
}

func (bc *SharedBattleComponent) EndBattle() {
	// 1. Send message with battle ended update
	bc.AddLogEntry("Battle ended!!!")
	bc.SendUpdate()
	// 2. Dsitribute loot to alive players if they won
	for _, p := range bc.Players {
		health_C := p.GetComponent("HealthComponent").(*HealthComponent)
		if health_C.IsAlive() {
			for _, e := range bc.Enemies {
				npc_C, npc_C_ok := e.GetComponent("NPCComponent").(*NPCComponent)
				if npc_C_ok {
					var drop []items.ItemBundle
					for _, cfg := range npc_C.PossibleLoot {
						chance := rand.Float64()
						if chance <= cfg.Chance {
							drop = append(drop, cfg.PossibleLoot)
						}
					}
					/* Creates a LootDispenserComponent and adds loot to be dropped to it. On the next update LootDispenser system
					will iterate over all enitites with this component and distribute loot */
					if len(drop) > 0 {
						lootDispenser_C := &LootDispenserComponent{
							State: ADDED,
						}
						lootDispenser_C.AddItems(drop...)
						p.AddComponent(lootDispenser_C)
					}
				}
			}
		}
	}

	// 3. Delete all NPC entities
	// 4. Remove battleComponents from players

}

func (bc *SharedBattleComponent) IsAnyoneStillLooting() bool {
	for _, p := range bc.Players {
		if p.HasComponent("LootDispenserComponent") {
			return true
		}
	}
	return false
}

func (bc *SharedBattleComponent) AddEnemy(e *entity.Entity) {
	bc.Enemies = append(bc.Enemies, e)
}

func (bc *SharedBattleComponent) AddPlayer(e *entity.Entity) {
	bc.Players = append(bc.Players, e)
}

/* Generates a status text of the battle that is used to display in the battle message,
text includes names and health of all entities in the battle */
func (bc SharedBattleComponent) generateStatusText() string {
	var playerStatus string
	for _, p := range bc.Players {
		name_C := p.GetComponent("NameComponent").(*NameComponent)
		health_c := p.GetComponent("HealthComponent").(*HealthComponent)
		playerStatus = fmt.Sprintf("%s\nHealth: %d/%d\n", name_C.GetName(), health_c.GetHealth(), health_c.GetMaxHealth())
	}
	status := playerStatus
	for _, e := range bc.Enemies {
		name_C := e.GetComponent("NameComponent").(*NameComponent)
		health_c := e.GetComponent("HealthComponent").(*HealthComponent)
		playerStatus = fmt.Sprintf("\n%s\nHealth: %d/%d\n", name_C.GetName(), health_c.GetHealth(), health_c.GetMaxHealth())
	}
	status += (playerStatus + "\n")
	for _, s := range bc.BattleLog {
		status += (s + "\n")
	}
	return status
}
