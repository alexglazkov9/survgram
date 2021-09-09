package activities

import (
	"fmt"
	"log"

	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/items/loot"
	"github.com/alexglazkov9/survgram/misc"
	"github.com/alexglazkov9/survgram/resources"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	ACTIVITIES_OPTIONS_NUMBER = 3
)

type Expedition struct {
	Bot     *tgbotapi.BotAPI
	Host    *entity.Entity
	Message *tgbotapi.Message

	LootManager    *loot.LootManager
	IsComplete     bool
	IsReadyForNext bool
	Location       Location

	currentActivity IActivity
	activityOptions []IActivity
}

func NewExpedition(bot *tgbotapi.BotAPI, char *entity.Entity, loc Location) *Expedition {
	return &Expedition{Bot: bot, Host: char, IsReadyForNext: false, Location: loc}
}

func (a *Expedition) Update(dt float64) {
	//Start next event
	if a.IsReadyForNext {
		a.next()
	}

	if a.currentActivity != nil && a.currentActivity.IsComplete() {
		a.currentActivity = nil
		a.sendNextKeyboard()
	}

	if a.currentActivity != nil && !a.currentActivity.IsComplete() {
		a.currentActivity.Update(dt)
	}
}

func (a *Expedition) SetSelectedActivity(i int) {
	log.Printf("%d selected", i)
	a.currentActivity = a.activityOptions[i]
	a.activityOptions = nil
	a.currentActivity.AddPlayer(a.Host)
	a.currentActivity.Start()
}

func (a *Expedition) next() {
	a.IsReadyForNext = false
	tgkb := misc.TGInlineKeyboard{Columns: 1}
	for i := 0; i < ACTIVITIES_OPTIONS_NUMBER; i++ {
		act := a.Location.GetActivity(a.Bot, a.LootManager)
		a.activityOptions = append(a.activityOptions, act)
		cbData := misc.CallbackData{Action: misc.ACTIVITY_SELECTED, Payload: fmt.Sprint(i)}
		tgkb.AddButton(act.GetDescription(), cbData.JSON())
	}
	editMsg := tgbotapi.NewEditMessageText(a.Message.Chat.ID, a.Message.MessageID, resources.ACTIVITY_SELECT_TEXT)
	editMarkup := tgbotapi.NewEditMessageReplyMarkup(a.Message.Chat.ID, a.Message.MessageID, tgkb.Generate())

	a.Bot.Send(editMsg)
	a.Bot.Send(editMarkup)
}

func (a *Expedition) sendNextKeyboard() {
	reply_markup := tgbotapi.NewInlineKeyboardMarkup()
	var row []tgbotapi.InlineKeyboardButton
	row = append(row, tgbotapi.NewInlineKeyboardButtonData(resources.EXPEDITION_CONTINUE, misc.CallbackData{Action: misc.EXPEDITION_CONTINUE}.JSON()))
	row = append(row, tgbotapi.NewInlineKeyboardButtonData(resources.EXPEDITION_LEAVE, misc.CallbackData{Action: misc.EXPEDITION_LEAVE}.JSON()))
	reply_markup.InlineKeyboard = append(reply_markup.InlineKeyboard, row)

	player_C := a.Host.GetComponent("PlayerComponent").(*components.PlayerComponent)
	msg := tgbotapi.NewMessage(player_C.ChatID, resources.EXPEDITION_MESSAGE_TEXT)
	msg.ReplyMarkup = reply_markup
	message, _ := a.Bot.Send(msg)
	a.Message = &message
}