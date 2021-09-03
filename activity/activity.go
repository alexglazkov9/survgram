package activity

import (
	"github.com/alexglazkov9/survgram/battle"
	"github.com/alexglazkov9/survgram/components"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/items"
	"github.com/alexglazkov9/survgram/lootmanager"
	"github.com/alexglazkov9/survgram/misc"
	"github.com/alexglazkov9/survgram/resources"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Activity struct {
	Bot            *tgbotapi.BotAPI
	Host           *entity.Entity
	Message        *tgbotapi.Message
	LootManager    *lootmanager.LootManager
	IsComplete     bool
	IsReadyForNext bool

	currentBattle *battle.Battle
}

func New(bot *tgbotapi.BotAPI, char *entity.Entity) *Activity {
	return &Activity{Bot: bot, Host: char, IsReadyForNext: true}
}

func (a *Activity) Update(dt float64) {
	//Start next event
	if a.IsReadyForNext {
		a.next()
	}

	if a.currentBattle != nil {
		if !a.currentBattle.Ended {
			a.currentBattle.Update(dt)
		} else {
			player_C := a.Host.GetComponent("PlayerComponent").(*components.PlayerComponent)
			if !a.LootManager.IsPlayerLooting(player_C.TelegramID) {
				a.currentBattle = nil
				a.sendNextKeyboard()
			}
		}
	}
}

func (a *Activity) next() {
	a.IsReadyForNext = false
	a.startBattle()
}

func (a *Activity) startBattle() {
	//Generate enemy
	enem := entity.New()
	enem.AddComponent(&components.NameComponent{
		Name: "Sobaka",
	})
	enem.AddComponent(&components.AttackComponent{
		AttackDamage: 1,
	})
	enem.AddComponent(&components.HealthComponent{
		MaxHealthPoints: 150,
		HealthPoints:    5,
	})
	enem.AddComponent(&components.EffectsComponent{})
	enem.AddComponent(&components.NPCComponent{
		PossibleLoot: []items.IItem{items.BaseItem{
			Name:   "Coin",
			Type:   "currency",
			Tier:   1,
			Rarity: items.COMMON,
		}},
	})

	//Setup battle
	a.currentBattle = &battle.Battle{Bot: a.Bot, LootManager: a.LootManager}
	a.currentBattle.AddToEnemies(enem)
	a.currentBattle.AddToPlayers(a.Host)
	a.currentBattle.Start()
}

func (a *Activity) sendNextKeyboard() {
	reply_markup := tgbotapi.NewInlineKeyboardMarkup()
	var row []tgbotapi.InlineKeyboardButton
	row = append(row, tgbotapi.NewInlineKeyboardButtonData(resources.ACTIVITY_CONTINUE, misc.CallbackData{Action: misc.ACTIVITY_CONTINUE}.JSON()))
	row = append(row, tgbotapi.NewInlineKeyboardButtonData(resources.ACTIVITY_LEAVE, misc.CallbackData{Action: misc.ACTIVITY_LEAVE}.JSON()))
	reply_markup.InlineKeyboard = append(reply_markup.InlineKeyboard, row)

	player_C := a.Host.GetComponent("PlayerComponent").(*components.PlayerComponent)
	msg := tgbotapi.NewMessage(player_C.ChatID, resources.ACTIVITY_MESSAGE_TEXT)
	msg.ReplyMarkup = reply_markup
	message, _ := a.Bot.Send(msg)
	a.Message = &message
}
