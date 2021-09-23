package components

import (
	"fmt"
	"log"

	"github.com/alexglazkov9/survgram/bot"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/items"
	"github.com/alexglazkov9/survgram/misc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type LootDispenserState int

const (
	ADDED   LootDispenserState = 1
	PICKING LootDispenserState = 2
)

type LootDispenserComponent struct {
	BaseComponent `bson:"-" json:"-"`

	LootDrop []items.ItemBundle
	State    LootDispenserState
	Message  *tgbotapi.Message

	SelectedItem *items.ItemBundle
}

func (ldc *LootDispenserComponent) Clone() entity.IComponent {
	copy := *ldc
	return &copy
}

func (ldc *LootDispenserComponent) AddItems(lootDrop ...items.ItemBundle) {
	ldc.LootDrop = append(ldc.LootDrop, lootDrop...)
}

func (ldc LootDispenserComponent) GenerateInlineKeyboard() *tgbotapi.InlineKeyboardMarkup {
	tg_kb := misc.TGInlineKeyboard{Columns: 2}

	//Header
	if ldc.SelectedItem != nil {
		cbData := misc.CallbackData{Action: misc.PICK_UP_ITEM, ID: fmt.Sprint(ldc.Parent.ID), Payload: fmt.Sprint(ldc.SelectedItem.ID)}
		tg_kb.AddHeaderButton("🔼Pick Up Selected", cbData.JSON())
	}
	cbData := misc.CallbackData{Action: misc.PICK_UP_ALL_ITEMS, ID: fmt.Sprint(ldc.Parent.ID)}
	tg_kb.AddHeaderButton("🔼Pick Up All", cbData.JSON())
	//Loot
	for _, bundle := range ldc.LootDrop {
		log.Println("HERERE!!!!")
		log.Println(ldc.Parent)
		log.Println(bundle)
		log.Println(bundle.GetItem().GetID())
		cbData := misc.CallbackData{Action: misc.SELECT_LOOT_ITEM, ID: fmt.Sprint(ldc.Parent.ID), Payload: fmt.Sprint(bundle.GetItem().GetID())}
		tg_kb.AddButton(fmt.Sprintf("%s (%d)", bundle.GetItem().GetName(), bundle.Qty), cbData.JSON())
	}
	//Dismiss
	cbData = misc.CallbackData{Action: misc.DISMISS_LOOT, ID: fmt.Sprint(ldc.Parent.ID)}
	tg_kb.AddFooterButton("Dismiss", cbData.JSON())

	return tg_kb.Generate()
}

func (ldc *LootDispenserComponent) Expire() {
	player_C := ldc.Parent.GetComponent("PlayerComponent").(*PlayerComponent)
	msg := tgbotapi.NewEditMessageText(player_C.ChatID, ldc.Message.MessageID, "Loot expired")
	bot.GetInstance().GetBot().Send(msg)
	ldc.Parent.RemoveComponent("LootDispenserComponent")
}

func (ldc *LootDispenserComponent) UpdateMessage() {
	if ldc.Message == nil {
		return
	}

	if len(ldc.LootDrop) == 0 {
		ldc.Expire()
		return
	}

	var text string
	if ldc.SelectedItem != nil {
		itm := items.GetItemCollection().GetItemById(ldc.SelectedItem.ID)
		text = fmt.Sprintf("Loot\n%s\n%s", itm.GetName(), itm.GetType())
	} else {
		text = "Loot"
	}

	player_C := ldc.Parent.GetComponent("PlayerComponent").(*PlayerComponent)
	editText := tgbotapi.NewEditMessageText(player_C.ChatID, ldc.Message.MessageID, text)
	editReplyMarkup := tgbotapi.NewEditMessageReplyMarkup(player_C.ChatID, ldc.Message.MessageID, *ldc.GenerateInlineKeyboard())

	_, err := bot.GetInstance().GetBot().Send(editText)
	if err != nil {
		log.Println(err)
	}
	_, err = bot.GetInstance().GetBot().Send(editReplyMarkup)
	if err != nil {
		log.Println(err)
	}
}
