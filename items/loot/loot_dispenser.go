package loot

import (
	"log"

	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/items"
	"github.com/alexglazkov9/survgram/misc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	EXPIRATION_TIME float64 = 15
)

type LootDispenser struct {
	Bot         tgbotapi.BotAPI
	pendingLoot map[int]*Loot
}

func NewLootManager(bot tgbotapi.BotAPI) *LootDispenser {
	return &LootDispenser{Bot: bot, pendingLoot: make(map[int]*Loot)}
}

func (lm *LootDispenser) Update(dt float64) {
	for key, elem := range lm.pendingLoot {
		//Remove already expired elements
		if elem.IsExpired() {
			log.Println("Remove expired")
			lm.Remove(key)
			continue
		}
		elem.ExpirationTime += dt
		if elem.ExpirationTime >= EXPIRATION_TIME {
			lm.Remove(key)
		}
	}
}

func (lm *LootDispenser) Add(target *entity.Entity, items []items.IItem) {
	//Remove existing
	player_C := target.GetComponent("PlayerComponent").(*components.PlayerComponent)
	lm.Remove(player_C.TelegramID)

	//Add new
	lm.pendingLoot[player_C.TelegramID] = NewLoot(lm.Bot, target, items)
	elem := lm.pendingLoot[player_C.TelegramID]
	elem.SendLoot()
}

func (lm *LootDispenser) Remove(key int) {
	if loot, ok := lm.pendingLoot[key]; ok {
		loot.Expire()
		delete(lm.pendingLoot, key)
	}
}

func (lm *LootDispenser) HandleInput(cbData misc.CallbackData, telegramIdFrom int) {
	switch cbData.Action {
	case misc.SELECT_LOOT_ITEM:
		loot := lm.pendingLoot[telegramIdFrom]
		loot.SetSelectedItem(cbData.Payload)
	case misc.PICK_UP_ITEM:
		loot := lm.pendingLoot[telegramIdFrom]
		loot.PickUp()
	case misc.PICK_UP_ALL_ITEMS:
	case misc.DISMISS_LOOT:
		lm.Remove(telegramIdFrom)
	}
}

/* Returns true if player hasn't finished
looting by picking all items or clicking DISMISS*/
func (lm *LootDispenser) IsPlayerLooting(telegram_id int) bool {
	if _, ok := lm.pendingLoot[telegram_id]; ok {
		return true
	}
	return false
}
