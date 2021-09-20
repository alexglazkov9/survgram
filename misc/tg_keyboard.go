package misc

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TGInlineKeyboard struct {
	headerButtons []tgbotapi.InlineKeyboardButton
	buttons       []tgbotapi.InlineKeyboardButton
	footerButtons []tgbotapi.InlineKeyboardButton
	Columns       int
}

func (tk TGInlineKeyboard) Generate() *tgbotapi.InlineKeyboardMarkup {
	kb_markup := tgbotapi.NewInlineKeyboardMarkup()
	kb_markup.InlineKeyboard = make([][]tgbotapi.InlineKeyboardButton, 0)
	var row []tgbotapi.InlineKeyboardButton
	i := 0

	//Setup header buttons
	if len(tk.headerButtons) > 0 {
		kb_markup.InlineKeyboard = append(kb_markup.InlineKeyboard, tk.headerButtons)
	}

	//Split buttons into columns
	for _, btn := range tk.buttons {
		row = append(row, btn)
		i++
		if i == tk.Columns {
			kb_markup.InlineKeyboard = append(kb_markup.InlineKeyboard, row)
			row = nil
			i = 0
		}
	}
	//Add the rest of buttons
	if len(row) > 0 {
		kb_markup.InlineKeyboard = append(kb_markup.InlineKeyboard, row)
	}

	//Setup footer buttons
	if len(tk.footerButtons) > 0 {
		kb_markup.InlineKeyboard = append(kb_markup.InlineKeyboard, tk.footerButtons)
	}

	return &kb_markup
}

//Specify number of columns for keyboard
func (tk *TGInlineKeyboard) SetColNumber(n int) {
	tk.Columns = n
}

func (tk *TGInlineKeyboard) AddButton(text string, data string) {
	tk.buttons = append(tk.buttons, tgbotapi.NewInlineKeyboardButtonData(text, data))
}

func (tk *TGInlineKeyboard) AddHeaderButton(text string, data string) {
	tk.headerButtons = append(tk.headerButtons, tgbotapi.NewInlineKeyboardButtonData(text, data))
}

func (tk *TGInlineKeyboard) AddFooterButton(text string, data string) {
	tk.footerButtons = append(tk.footerButtons, tgbotapi.NewInlineKeyboardButtonData(text, data))
}
