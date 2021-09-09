package game

import (
	"fmt"
	"log"
	"strconv"

	"github.com/alexglazkov9/survgram/activities"
	"github.com/alexglazkov9/survgram/database"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/items/loot"
	"github.com/alexglazkov9/survgram/misc"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	gl "github.com/kutase/go-gameloop"
)

type Game struct {
	Locations        *activities.Locations
	CharacterManager *database.CharacterManager
	Expeditions      *activities.Expeditions
	LootManager      *loot.LootManager
	Bot              *tgbotapi.BotAPI
	Engine           *gl.GameLoop
}

// GetInstance - returns instance of Character Manager
func New(bot *tgbotapi.BotAPI) *Game {
	instance := &Game{}
	instance.Bot = bot
	instance.Locations = activities.NewLocations()
	instance.CharacterManager = database.NewCharacterManager(database.GetInstance(), instance.Locations.GetStartLocation())
	instance.LootManager = loot.NewLootManager(*bot)
	instance.Expeditions = &activities.Expeditions{LootManager: instance.LootManager, CharacterManager: instance.CharacterManager}
	instance.Engine = gl.New(30, func(dt float64) {
		instance.Expeditions.Update(dt)
		instance.LootManager.Update(dt)
	})
	instance.Engine.Start()

	return instance
}

func (g Game) HandleInput(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		callbackData := misc.CallbackData{}
		callbackData.FromJSON(update.CallbackQuery.Data)

		chrctr := g.CharacterManager.GetCharacter(update.CallbackQuery.From.ID)
		switch update.CallbackQuery.Data {
		case "goto":
			chrctr := g.CharacterManager.GetCharacter(update.CallbackQuery.From.ID)
			player_C := chrctr.GetComponent("PlayerComponent").(*components.PlayerComponent)

			var buttons tgbotapi.InlineKeyboardMarkup

			loc := g.Locations.GetLocation(player_C.CurrentLocation)

			//Add destinations to the keyboard
			var row []tgbotapi.InlineKeyboardButton
			i := 0
			for _, dest := range loc.Destinations {
				cbData := misc.CallbackData{Action: misc.GO_TO, Payload: fmt.Sprint(dest.GetID())}
				row = append(row, tgbotapi.NewInlineKeyboardButtonData(dest.Name, cbData.JSON()))
				i++
				//Change number of columns
				if i == 2 {
					buttons.InlineKeyboard = append(buttons.InlineKeyboard, row)
					row = nil
					i = 0
				}
			}
			//Add the rest of the destinations
			if len(row) > 0 {
				buttons.InlineKeyboard = append(buttons.InlineKeyboard, row)
			}

			textEdit := tgbotapi.NewEditMessageText(player_C.ChatID, update.CallbackQuery.Message.MessageID, "Go to ...")
			markupEdit := tgbotapi.NewEditMessageReplyMarkup(player_C.ChatID, update.CallbackQuery.Message.MessageID, buttons)
			g.Bot.Send(textEdit)
			g.Bot.Send(markupEdit)
		case "do":
			g.Expeditions.Add(activities.NewExpedition(g.Bot, chrctr, *activities.NewLocation(1, "Test location")))
		}

		//Moves character to a new location
		switch callbackData.Action {
		case misc.GO_TO:
			chrctr := g.CharacterManager.GetCharacter(update.CallbackQuery.From.ID)
			player_C := chrctr.GetComponent("PlayerComponent").(*components.PlayerComponent)

			objID, err := strconv.Atoi(callbackData.Payload)
			if err != nil {
				log.Fatal(err)
			}
			player_C.CurrentLocation = objID
			msg := tgbotapi.NewMessage(player_C.ChatID, "You have reached the new location")
			g.Bot.Send(msg)
		case misc.SELECT_LOOT_ITEM, misc.PICK_UP_ITEM, misc.PICK_UP_ALL_ITEMS, misc.DISMISS_LOOT:
			g.LootManager.HandleInput(callbackData, update.CallbackQuery.From.ID)
		case misc.EXPEDITION_CONTINUE, misc.EXPEDITION_LEAVE, misc.ACTIVITY_SELECTED:
			g.Expeditions.HandleInput(update)
		}
	}
}
