package systems

import (
	"fmt"
	"math/rand"

	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/combat"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const UPDATE_PERIOD float64 = 2

type BattleSystem struct {
	Bot       *tgbotapi.BotAPI
	currentDt float64

	manager         *entity.Manager
	characterHelper interfaces.CharacterHelper
}

func NewBattleSystem(manager *entity.Manager, characterHelper interfaces.CharacterHelper) *BattleSystem {
	battle := &BattleSystem{manager: manager, characterHelper: characterHelper}

	return battle
}

func (b *BattleSystem) Update(dt float64) {
	b.currentDt = dt

	for _, entity := range b.manager.QueryEntities("SharedBattleComponent") {
		battle_C := entity.GetComponent("SharedBattleComponent").(*components.SharedBattleComponent)

		switch battle_C.CurrentState {
		case components.PREACTIVITY:
			battle_C.StartBattle()
		case components.RUNNING:
			for _, e := range battle_C.GetAllEntities(true) {
				b.handleAttacks(e, battle_C)
			}
			b.handleMessageUpdate(entity)
		case components.POSTACTIVITY:
			// Check if anyone still looting and finish battle
			if !battle_C.IsAnyoneStillLooting() {
				entity.GetComponent("ActivityStatusComponent").(*components.ActivityStatusComponent).IsComplete = true
				for _, p := range battle_C.Players {
					b.characterHelper.UpdateCharacter(p)
				}
			}
		}
	}
}

/* Handles attacks and abilities from all entities in the battle */
func (b *BattleSystem) handleAttacks(entity *entity.Entity, battle_C *components.SharedBattleComponent) {
	attack_C := entity.GetComponent("AttackComponent").(*components.AttackComponent)
	ability_C, ability_C_ok := entity.GetComponent("AbilityComponent").(*components.AbilityComponent)
	name_C := entity.GetComponent("NameComponent").(*components.NameComponent)

	attack_C.NextAttackTimer += b.currentDt

	//Cast ability if ready
	if ability_C_ok && ability_C.IsReadyToCast() {
		enms := battle_C.GetEnemiesForEntity(entity, true)
		target := enms[rand.Intn(len(enms))]
		target.GetComponent("EffectsComponent").(*components.EffectsComponent).ApplyEffect(ability_C.Ability.Effects)
		ability_C.Energy -= ability_C.Ability.Energy_cost

		attack_C.NextAttackTimer = 0 //reset attack timer
	}

	//Attack if reaedy
	if attack_C.IsReadyToAttack() {
		attack_C.NextAttackTimer -= components.ATTACK_SPEED //reset attack timer

		var attackDetails combat.AttackDetails
		enms := battle_C.GetEnemiesForEntity(entity, true)
		target := enms[rand.Intn(len(enms))]

		attackDetails.Damage = target.GetComponent("HealthComponent").(*components.HealthComponent).ApplyDamage(attack_C.GetBasicAttackDamage())
		attackDetails.AttackType = combat.PhysicalAttack
		attackDetails.Target = target.GetComponent("NameComponent").(*components.NameComponent).GetName()
		attackDetails.Attacker = name_C.GetName()
		if ability_C_ok {
			ability_C.GainEnergy(attackDetails.Damage)
		}
		battle_C.AddLogEntry(fmt.Sprintf("%s ⚔ %s - %d damage", attackDetails.Attacker, attackDetails.Target, attackDetails.Damage))
	}

	//End battle if all entities on either side are dead
	if battle_C.HasAnyoneWon() {
		battle_C.CurrentState = components.POSTACTIVITY
		battle_C.EndBattle()
		//Dispose of all non-player entities
		for _, e := range append(battle_C.Players, battle_C.Enemies...) {
			if !e.HasComponent("PlayerComponent") {
				b.manager.DisposeEntity(e)
			}
		}
		// 4. Remove battleComponents from players
	}
}

/* Updates the message for all players with battlelog info */
func (b *BattleSystem) handleMessageUpdate(e *entity.Entity) {
	battle_C := e.GetComponent("SharedBattleComponent").(*components.SharedBattleComponent)
	battle_C.NextUpdateTimer += b.currentDt
	if battle_C.NextUpdateTimer >= UPDATE_PERIOD {
		battle_C.NextUpdateTimer -= UPDATE_PERIOD
		battle_C.SendUpdate()
	}
}
