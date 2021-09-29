package components

import (
	"fmt"
	"math/rand"

	"github.com/alexglazkov9/survgram/bot"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/items"
	"github.com/alexglazkov9/survgram/misc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type GatheringActivityComponent struct {
	BaseComponent `bson:"-" json:"-"`

	IsSuccesful  bool
	Players      []*entity.Entity
	Messages     map[int]tgbotapi.Message //Holds tgbotapi.Message sent to players in map[TelegramID_of_the_player]players_message in the current activity
	CurrentState ActivityState
	Resource     items.ItemBundle
	Count        int
}

/*
	Generates message text according to current activity status and adds a keyboard if activity still active.
	Checks if message was provided - then updates it, otherwise sends a new one
*/
func (ga *GatheringActivityComponent) SendUpdate(e *entity.Entity) {
	player_C := e.GetComponent("PlayerComponent").(*PlayerComponent)
	var msg tgbotapi.Chattable

	//Default text for update
	msg_text := fmt.Sprintf("Pick the correct button to gather %d/%d", ga.Count, 4)

	//End activity text
	if ga.CurrentState == POSTACTIVITY {
		if ga.IsSuccesful {
			msg_text = fmt.Sprintf("You succesfully gathered %s", ga.Resource.GetItem().GetName())
		} else {
			name_C := e.GetComponent("NameComponent").(*NameComponent)
			msg_text = fmt.Sprintf("Gathering failed by %s", name_C.GetName())
		}
	}

	if _, ok := ga.Messages[player_C.TelegramID]; ok {
		//Updates the message
		temp := tgbotapi.NewEditMessageText(
			player_C.ChatID,
			ga.Messages[player_C.TelegramID].MessageID,
			msg_text,
		)
		//Adds keyboard only if activity is active
		if ga.CurrentState != POSTACTIVITY {
			kb := misc.TGInlineKeyboard{Columns: 2}
			btn_index := rand.Intn(4)
			for i := 0; i < 4; i++ {
				if btn_index == i {
					kb.AddButton(
						"Gather",
						misc.CallbackData{
							Action: misc.GATHERING_CORRECT,
							ID:     fmt.Sprint(ga.Parent.ID)}.JSON(),
					)
				} else {
					kb.AddButton(
						"-",
						misc.CallbackData{
							Action: misc.GATHERING_INCORRECT,
							ID:     fmt.Sprint(ga.Parent.ID)}.JSON(),
					)
				}
			}

			reply_markup := kb.Generate()
			temp.ReplyMarkup = reply_markup
		}
		msg = temp
	} else {
		//Sends a new message
		temp := tgbotapi.NewMessage(
			player_C.ChatID,
			msg_text,
		)
		//Adds keyboard only if activity is active
		if ga.CurrentState != POSTACTIVITY {
			kb := misc.TGInlineKeyboard{Columns: 2}
			btn_index := rand.Intn(4)
			for i := 0; i < 4; i++ {
				if btn_index == i {
					kb.AddButton(
						"Gather",
						misc.CallbackData{
							Action: misc.GATHERING_CORRECT,
							ID:     fmt.Sprint(ga.Parent.ID)}.JSON(),
					)
				} else {
					kb.AddButton(
						"-",
						misc.CallbackData{
							Action: misc.GATHERING_INCORRECT,
							ID:     fmt.Sprint(ga.Parent.ID)}.JSON(),
					)
				}
			}

			reply_markup := *kb.Generate()
			temp.ReplyMarkup = reply_markup
		}
		msg = temp
	}
	ga.Messages[player_C.TelegramID], _ = bot.GetInstance().GetBot().Send(msg)
}
