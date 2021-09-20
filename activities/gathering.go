package activities

import (
	"fmt"
	"math/rand"

	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/items"
	"github.com/alexglazkov9/survgram/items/loot"
	"github.com/alexglazkov9/survgram/misc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type GatheringActivity struct {
	Bot                *tgbotapi.BotAPI
	IsActivityComplete bool
	LootDispenser      *loot.LootDispenser
	SelectedResource   items.IItem

	WrongActionPicked bool

	/* Holds tgbotapi.Message sent to players in map[TelegramID_of_the_player]players_message in the current activity */
	messages     map[int]tgbotapi.Message
	players      []*entity.Entity
	currentState ActivityState

	//Count succesfull gatherings
	count int
}

func NewGatheringActivity(bot *tgbotapi.BotAPI, lm *loot.LootDispenser, spawn_chances []SpawnChance) *GatheringActivity {
	messages := make(map[int]tgbotapi.Message)
	itm_id := spawn_chances[rand.Intn(len(spawn_chances))].Id
	selected_resource := items.GetItemCollection().GetItemById(itm_id)
	actvty := &GatheringActivity{Bot: bot, LootDispenser: lm, messages: messages, SelectedResource: selected_resource, currentState: PREACTIVITY}

	return actvty
}

func (ga *GatheringActivity) Start() {
	for _, p := range ga.players {
		player_C := p.GetComponent("PlayerComponent").(*components.PlayerComponent)
		if _, ok := ga.messages[player_C.TelegramID]; ok {
			msg := tgbotapi.NewEditMessageText(player_C.ChatID, ga.messages[player_C.TelegramID].MessageID, ga.generateStatusText())
			reply_markup := tgbotapi.NewEditMessageReplyMarkup(player_C.ChatID, ga.messages[player_C.TelegramID].MessageID, ga.generateReplyMarkup())
			ga.messages[player_C.TelegramID], _ = ga.Bot.Send(msg)
			ga.Bot.Send(reply_markup)
		} else {
			msg := tgbotapi.NewMessage(player_C.ChatID, ga.generateStatusText())
			msg.ReplyMarkup = ga.generateReplyMarkup()
			ga.messages[player_C.TelegramID], _ = ga.Bot.Send(msg)
		}
	}
	ga.currentState = RUNNING
}

func (ga *GatheringActivity) Update(dt float64) {
	switch ga.currentState {
	case PREACTIVITY:
	case RUNNING:
		if ga.WrongActionPicked || ga.count >= 4 {
			ga.SetState(POSTACTIVITY)
		}
	case POSTACTIVITY:
		still_looting := false
		for _, p := range ga.players {
			player_C := p.GetComponent("PlayerComponent").(*components.PlayerComponent)
			if ga.LootDispenser.IsPlayerLooting(player_C.TelegramID) {
				still_looting = true
			}
		}
		if !still_looting {
			ga.IsActivityComplete = true
		}
	}
}
func (ga *GatheringActivity) IsComplete() bool { return ga.IsActivityComplete }
func (ga *GatheringActivity) GetDescription() string {
	return fmt.Sprintf("⛏ %s", ga.SelectedResource.GetName())
}
func (ga *GatheringActivity) AddPlayer(player *entity.Entity, msg ...*tgbotapi.Message) {
	ga.players = append(ga.players, player)
	if len(msg) > 0 {
		player_C := player.GetComponent("PlayerComponent").(*components.PlayerComponent)
		ga.messages[player_C.TelegramID] = *msg[0]
	}
}

func (ga *GatheringActivity) generateStatusText() string {
	return "Gather smth"
}

func (ga *GatheringActivity) generateReplyMarkup() tgbotapi.InlineKeyboardMarkup {
	rm := misc.TGInlineKeyboard{Columns: 2}
	btn_index := rand.Intn(4)
	for i := 0; i < 4; i++ {
		if btn_index == i {
			rm.AddButton("Gather", misc.CallbackData{Action: misc.GATHERING_CORRECT}.JSON())
		} else {
			rm.AddButton("-", misc.CallbackData{Action: misc.GATHERING_INCORRECT}.JSON())
		}
	}

	return *rm.Generate()
}

func (ga *GatheringActivity) endActivity() {
	for _, p := range ga.players {
		player_C := p.GetComponent("PlayerComponent").(*components.PlayerComponent)
		text := ""
		if ga.WrongActionPicked {
			text = "Failed"
		} else {
			text = "Success"
			ga.LootDispenser.Add(p, []items.ItemBundle{{ga.SelectedResource.GetID(), 1}})
		}
		msg := tgbotapi.NewEditMessageText(player_C.ChatID, ga.messages[player_C.TelegramID].MessageID, text)

		ga.messages[player_C.TelegramID], _ = ga.Bot.Send(msg)
	}
}

func (ga *GatheringActivity) SendUpdate() {
	for _, p := range ga.players {
		player_C := p.GetComponent("PlayerComponent").(*components.PlayerComponent)
		msg := tgbotapi.NewEditMessageText(player_C.ChatID, ga.messages[player_C.TelegramID].MessageID, ga.generateStatusText())
		reply_markup := tgbotapi.NewEditMessageReplyMarkup(player_C.ChatID, ga.messages[player_C.TelegramID].MessageID, ga.generateReplyMarkup())
		ga.messages[player_C.TelegramID], _ = ga.Bot.Send(msg)
		ga.Bot.Send(reply_markup)
		ga.count++
	}
}

func (ga *GatheringActivity) SetState(state ActivityState) {
	switch state {
	case POSTACTIVITY:
		ga.endActivity()
		ga.currentState = state
	}
}
