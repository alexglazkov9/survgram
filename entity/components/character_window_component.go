package components

import (
	"fmt"
	"log"

	"github.com/alexglazkov9/survgram/bot"
	"github.com/alexglazkov9/survgram/misc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type CharacterWindowTab int

const (
	CHARACTER CharacterWindowTab = iota
	EQUIPMENT
)

type CharacterWindowComponent struct {
	BaseComponent

	IsSent       bool
	CurrentTab   CharacterWindowTab
	SelectedItem *ItemID
	Message      *tgbotapi.Message
}

func (cwc *CharacterWindowComponent) SendUpdate() {
	name_C := cwc.Parent.GetComponent("NameComponent").(*NameComponent)
	player_C := cwc.Parent.GetComponent("PlayerComponent").(*PlayerComponent)
	stats_C := cwc.Parent.GetComponent("StatsComponent").(*StatsComponent)
	health_C := cwc.Parent.GetComponent("HealthComponent").(*HealthComponent)
	equipment_C := cwc.Parent.GetComponent("PlayerEquipmentComponent").(*PlayerEquipmentComponent)

	var text string
	tgkb := &misc.TGInlineKeyboard{Columns: 2, IsClosable: true}

	//Create tab switching buttons
	cb_data := misc.CallbackData{Action: misc.CHARACTER_SELECT_TAB, Payload: fmt.Sprint(CHARACTER)}
	tgkb.AddHeaderButton("Character", cb_data.JSON())

	cb_data = misc.CallbackData{Action: misc.CHARACTER_SELECT_TAB, Payload: fmt.Sprint(EQUIPMENT)}
	tgkb.AddHeaderButton("Equipment", cb_data.JSON())
	log.Println("about to send")
	switch cwc.CurrentTab {
	case CHARACTER:
		text = fmt.Sprintf("%s\n\nHP:%d\\%d\nSTR:%d\nAGI:%d\nINT:%d", name_C.GetName(), health_C.HealthPoints, health_C.MaxHealthPoints, stats_C.Str, stats_C.Agi, stats_C.Int)
	case EQUIPMENT:
		text = "Equipment"
		//Prepare buttons for each piece of equipment
		log.Println(equipment_C.GetAll())
		for _, item := range equipment_C.GetAll() {
			//Passes an itemId as parameter
			if item != nil {
				log.Println("not nill")
				tgkb.AddButton(item.GetItem().GetName(), cb_data.JSON())
				cb_data = misc.CallbackData{Action: misc.CHARACTER_EQUIPMENT_SELECTED, Payload: fmt.Sprint(item.GetItem().GetID())}
			} else {
				tgkb.AddButton("-", "-")
			}
		}
		if cwc.SelectedItem != nil {
			text += fmt.Sprintf("\n\n%s", cwc.SelectedItem.GetItem().GetFormattedItem(true))
		}
	}

	var msg tgbotapi.Chattable
	if cwc.Message != nil {
		msg_t := tgbotapi.NewEditMessageCaption(player_C.ChatID, cwc.Message.MessageID, text)
		msg_t.ReplyMarkup = tgkb.Generate()
		msg_t.ParseMode = tgbotapi.ModeMarkdown
		msg = msg_t
	} else {
		//Prepare message
		msg_t := tgbotapi.NewPhotoUpload(player_C.ChatID, "./resources/images/character.jpg")
		msg_t.Caption = text
		msg_t.ReplyMarkup = tgkb.Generate()
		msg_t.ParseMode = tgbotapi.ModeMarkdown
		msg = msg_t
	}
	log.Println("sending")
	sent_msg, err := bot.GetInstance().GetBot().Send(msg)
	if err != nil {
		log.Println(err)
	}
	cwc.Message = &sent_msg
}
