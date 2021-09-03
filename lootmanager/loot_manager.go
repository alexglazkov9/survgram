package lootmanager

import (
	"log"

	"github.com/alexglazkov9/survgram/components"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/items"
	"github.com/alexglazkov9/survgram/loot"
	"github.com/alexglazkov9/survgram/misc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	EXPIRATION_TIME float64 = 5
)

type LootManager struct {
	Bot         tgbotapi.BotAPI
	pendingLoot map[int]*loot.Loot
}

func New(bot tgbotapi.BotAPI) *LootManager {
	return &LootManager{Bot: bot, pendingLoot: make(map[int]*loot.Loot)}
}

func (lm *LootManager) Update(dt float64) {
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

func (lm *LootManager) Add(target *entity.Entity, items []items.IItem) {
	//Remove existing
	player_C := target.GetComponent("PlayerComponent").(*components.PlayerComponent)
	lm.Remove(player_C.TelegramID)

	//Add new
	lm.pendingLoot[player_C.TelegramID] = loot.New(lm.Bot, target, items)
	elem := lm.pendingLoot[player_C.TelegramID]
	elem.SendLoot()
}

func (lm *LootManager) Remove(key int) {
	if loot, ok := lm.pendingLoot[key]; ok {
		loot.Expire()
		delete(lm.pendingLoot, key)
	}
}

func (lm *LootManager) HandleInput(cbData misc.CallbackData, telegramIdFrom int) {
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

func (lm *LootManager) IsPlayerLooting(telegram_id int) bool {
	if _, ok := lm.pendingLoot[telegram_id]; ok {
		return true
	}
	return false
}
