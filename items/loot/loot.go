package loot

import (
	"fmt"
	"log"

	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/items"
	"github.com/alexglazkov9/survgram/misc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Loot struct {
	Bot     tgbotapi.BotAPI
	Target  *entity.Entity
	Items   []items.ItemBundle
	message *tgbotapi.Message

	selectedItem   *items.ItemBundle
	ExpirationTime float64

	isExpired bool
}

func NewLoot(bot tgbotapi.BotAPI, target *entity.Entity, items []items.ItemBundle) *Loot {
	return &Loot{
		Bot:          bot,
		Target:       target,
		Items:        items,
		selectedItem: nil,
	}
}

func (l *Loot) SendLoot() {
	text := "Loot"

	player_C := l.Target.GetComponent("PlayerComponent").(*components.PlayerComponent)
	msg := tgbotapi.NewMessage(player_C.ChatID, text)
	msg.ReplyMarkup = l.generateInlineKeyboard()
	message, _ := l.Bot.Send(msg)
	l.message = &message
}

func (l *Loot) UpdateMessage() {
	if l.message == nil {
		return
	}

	if len(l.Items) == 0 {
		l.Expire()
		return
	}

	var text string
	if l.selectedItem != nil {
		itm := items.GetItemCollection().GetItemById(l.selectedItem.ID)
		text = fmt.Sprintf("Loot\n%s\n%s", itm.GetName(), itm.GetType())
	} else {
		text = "Loot"
	}

	player_C := l.Target.GetComponent("PlayerComponent").(*components.PlayerComponent)
	editText := tgbotapi.NewEditMessageText(player_C.ChatID, l.message.MessageID, text)
	editReplyMarkup := tgbotapi.NewEditMessageReplyMarkup(player_C.ChatID, l.message.MessageID, l.generateInlineKeyboard())

	_, err := l.Bot.Send(editText)
	if err != nil {
		log.Println(err)
	}
	_, err = l.Bot.Send(editReplyMarkup)
	if err != nil {
		log.Println(err)
	}
}

func (l *Loot) Expire() {
	player_C := l.Target.GetComponent("PlayerComponent").(*components.PlayerComponent)
	msg := tgbotapi.NewEditMessageText(player_C.ChatID, l.message.MessageID, "Loot expired")
	l.Bot.Send(msg)
	l.isExpired = true
}

func (l *Loot) SetSelectedItem(id int) {
	for _, item := range l.Items {
		if item.ID == id {
			l.selectedItem = &item
			l.UpdateMessage()
			return
		}
	}

}

func (l *Loot) PickUp() {
	l.Target.GetComponent("InventoryComponent").(*components.InventoryComponent).AddItems(*l.selectedItem)
	for i, item := range l.Items {
		if item.ID == l.selectedItem.ID {
			l.Items[i] = l.Items[len(l.Items)-1]
			l.Items = l.Items[:len(l.Items)-1]
			l.selectedItem = nil
			l.UpdateMessage()

			return
		}
	}
}

func (l Loot) IsExpired() bool {
	return l.isExpired
}

func (l Loot) generateInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	tg_kb := misc.TGInlineKeyboard{}

	//Header
	if l.selectedItem != nil {
		cbData := misc.CallbackData{Action: misc.PICK_UP_ITEM, Payload: fmt.Sprint(l.selectedItem.ID)}
		tg_kb.AddHeaderButton("Pick Up Selected", cbData.JSON())
	}
	cbData := misc.CallbackData{Action: misc.PICK_UP_ALL_ITEMS}
	tg_kb.AddHeaderButton("Pick Up All", cbData.JSON())
	//Loot
	for _, item_id := range l.Items {
		itm := items.GetItemCollection().GetItemById(item_id.ID)
		cbData := misc.CallbackData{Action: misc.SELECT_LOOT_ITEM, Payload: fmt.Sprint(itm.GetID())}
		tg_kb.AddButton(itm.GetName(), cbData.JSON())
	}
	//Dismiss
	cbData = misc.CallbackData{Action: misc.DISMISS_LOOT}
	tg_kb.AddFooterButton("Dismiss", cbData.JSON())

	return tg_kb.Generate()
}
