package activities

import (
	"fmt"
	"log"

	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/entity/enemies"
	"github.com/alexglazkov9/survgram/items"
	"github.com/alexglazkov9/survgram/items/loot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Side int

const (
	PLAYERS Side = iota
	ENEMIES
)

const UPDATE_PERIOD float64 = 2

type Battle struct {
	Bot              *tgbotapi.BotAPI
	IsBattleComplete bool
	enemies          []*entity.Entity
	players          []*entity.Entity
	messages         map[int]tgbotapi.Message //Holds tgbotapi.Message sent to players in map[TelegramID_of_the_player]players_message in the current activity
	currentState     ActivityState
	battleLog        []string
	nextUpdateTimer  float64
	lootDispenser    *loot.LootDispenser
}

func NewBattle(bot *tgbotapi.BotAPI, lm *loot.LootDispenser, spawn_chances []SpawnChance) *Battle {
	enem := enemies.GetInstance().GetEnemyById(GetSpawneeId(spawn_chances))
	enem.AddComponent(&components.EffectsComponent{})
	possibleLoot := make([]items.IItem, 0)
	possibleLoot = append(possibleLoot, items.GetInstance().GetItemById(1))
	possibleLoot = append(possibleLoot, items.GetInstance().GetItemById(2))
	enem.AddComponent(&components.NPCComponent{
		PossibleLoot: possibleLoot,
	})
	messages := make(map[int]tgbotapi.Message)
	battle := &Battle{Bot: bot, lootDispenser: lm, messages: messages}
	battle.AddEnemy(&enem)
	battle.currentState = PREACTIVITY

	return battle
}

func (b *Battle) Update(dt float64) {
	for _, e := range b.getAllEntities() {
		e.Update(dt)
	}

	switch b.currentState {
	case PREACTIVITY:
	case RUNNING:
		b.handleAttacks()
		b.handleMessageUpdate(dt)
	case POSTACTIVITY:
		//If any of the players still looting do nothing
		anyone_looting := false
		for _, p := range b.getAllPlayers(true) {
			player_C := p.GetComponent("PlayerComponent").(*components.PlayerComponent)
			if b.lootDispenser.IsPlayerLooting(player_C.TelegramID) {
				anyone_looting = true
				break
			}
		}
		//Otherwise raise battle complete flag and parent class will remove it
		if !anyone_looting {
			b.IsBattleComplete = true
		}
	}
}

func (b *Battle) Start() {
	b.battleLog = make([]string, 0)
	for _, p := range b.players {
		player_C := p.GetComponent("PlayerComponent").(*components.PlayerComponent)
		var msg tgbotapi.Chattable
		if _, ok := b.messages[player_C.TelegramID]; ok {
			msg = tgbotapi.NewEditMessageText(player_C.ChatID, b.messages[player_C.TelegramID].MessageID, b.generateStatusText())
		} else {
			msg = tgbotapi.NewMessage(player_C.ChatID, b.generateStatusText())
		}
		b.messages[player_C.TelegramID], _ = b.Bot.Send(msg)
	}
	b.currentState = RUNNING
}

func (b *Battle) AddEnemy(e *entity.Entity) {
	e.AddComponent(&components.LogComponent{
		Log: &b.battleLog,
	})
	b.enemies = append(b.enemies, e)
}

func (b *Battle) AddPlayer(e *entity.Entity, msg ...*tgbotapi.Message) {
	e.AddComponent(&components.LogComponent{
		Log: &b.battleLog,
	})
	b.players = append(b.players, e)
	if len(msg) > 0 {
		player_C := e.GetComponent("PlayerComponent").(*components.PlayerComponent)
		b.messages[player_C.TelegramID] = *msg[0]
	}
}

func (b *Battle) IsComplete() bool {
	return b.IsBattleComplete
}

func (b *Battle) GetDescription() string {
	descr := ""
	for _, e := range b.enemies {
		name_C := e.GetComponent("NameComponent").(*components.NameComponent)
		descr += (name_C.GetName() + "\n")
	}
	return fmt.Sprintf("⚔ %s", descr)
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
		return len(b.getAllPlayers(true)) > 0
	}
	if side == ENEMIES {
		return len(b.getAllEnemies(true)) > 0
	}
	return false
}

/* Returns list of all players in the battle, pass true to return ALIVE players only */
func (b Battle) getAllPlayers(args ...bool) []*entity.Entity {
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

/* Returns list of all enemies in the battle, pass true to return ALIVE enemies only */
func (b Battle) getAllEnemies(args ...bool) []*entity.Entity {
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

	return b.enemies
}

func (b *Battle) endBattle() {
	b.currentState = POSTACTIVITY
	b.sendUpdate()
	//TODO Possibly need to remove any effects still applied on the character
	for _, e := range b.getAllEntities() {
		e.RemoveComponent("LogComponent")
	}
	for _, e := range b.getAllEnemies() {
		log.Println(e.Components)
		npc_C := e.GetComponent("NPCComponent").(*components.NPCComponent)
		for _, p := range b.getAllPlayers(true) {
			b.lootDispenser.Add(p, npc_C.PossibleLoot)
		}
	}
}

/* Handles attacks and abilities from all entities in the battle */
func (b *Battle) handleAttacks() {
	//Handle enemies
	for _, e := range b.getAllEnemies(true) {
		//Hadnle abilities
		ability_c, ok := e.GetComponent("AbilityComponent").(*components.AbilityComponent)
		if ok {
			if ability_c.IsReadyToCast() {
				ability_c.CastAbility(b.enemies, b.players)
			}
		}
		//Handle auto-attacks
		attack_c := e.GetComponent("AttackComponent").(*components.AttackComponent)
		if attack_c.IsReadyToAttack() {
			attackDetails := attack_c.Attack(b.players)
			attackDetails.Attacker = e.GetComponent("NameComponent").(*components.NameComponent).GetName()
			if ability_c != nil {
				ability_c.GainEnergy(attackDetails.Damage)
			}
			b.battleLog = append(b.battleLog, fmt.Sprintf("%s attacks %s for %d damage", attackDetails.Attacker, attackDetails.Target, attackDetails.Damage))
		}
	}

	//Handle players
	for _, p := range b.getAllPlayers(true) {
		//Hadnle abilities
		ability_c, ok := p.GetComponent("AbilityComponent").(*components.AbilityComponent)
		if ok {
			if ability_c.IsReadyToCast() {
				ability_c.CastAbility(b.players, b.enemies)
			}
		}
		//Handle auto-attacks
		attack_c := p.GetComponent("AttackComponent").(*components.AttackComponent)
		if attack_c.IsReadyToAttack() {
			attackDetails := attack_c.Attack(b.enemies)
			attackDetails.Attacker = p.GetComponent("NameComponent").(*components.NameComponent).GetName()
			if ability_c != nil {
				ability_c.GainEnergy(attackDetails.Damage)
			}
			b.battleLog = append(b.battleLog, fmt.Sprintf("%s attacks %s for %d damage", attackDetails.Attacker, attackDetails.Target, attackDetails.Damage))
		}
	}

	//End battle if all entities on either side are dead
	if !b.isAnyoneAlive(ENEMIES) || !b.isAnyoneAlive(PLAYERS) {
		b.endBattle()
	}
}

/* Updates the message for all players with battlelog info */
func (b *Battle) handleMessageUpdate(dt float64) {
	b.nextUpdateTimer += dt
	if b.nextUpdateTimer >= UPDATE_PERIOD {
		b.nextUpdateTimer = 0
		b.sendUpdate()
	}
}
