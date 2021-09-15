package game

import (
	"fmt"
	"log"
	"strconv"

	"github.com/alexglazkov9/survgram/activities"
	"github.com/alexglazkov9/survgram/database"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/entity/enemies"
	"github.com/alexglazkov9/survgram/items"
	"github.com/alexglazkov9/survgram/items/loot"
	"github.com/alexglazkov9/survgram/misc"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	gl "github.com/kutase/go-gameloop"
)

type Game struct {
	CharacterManager *database.CharacterManager
	Expeditions      *activities.Expeditions
	LootManager      *loot.LootDispenser
	Bot              *tgbotapi.BotAPI
	Engine           *gl.GameLoop
}

// GetInstance - returns instance of Character Manager
func New(bot *tgbotapi.BotAPI) *Game {
	instance := &Game{}
	instance.Bot = bot
	instance.CharacterManager = database.NewCharacterManager(database.GetInstance(), activities.GetLocations().GetStartLocation())
	instance.LootManager = loot.NewLootManager(*bot)
	instance.Expeditions = &activities.Expeditions{LootManager: instance.LootManager, CharacterManager: instance.CharacterManager}
	instance.Engine = gl.New(30, func(dt float64) {
		instance.Expeditions.Update(dt)
		instance.LootManager.Update(dt)
	})
	instance.Engine.Start()
	enemies.GetInstance()
	items.GetItemCollection()
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

			//Add destinations to the keyboard
			loc := activities.GetLocations().GetLocation(player_C.CurrentLocation)
			kb := misc.TGInlineKeyboard{Columns: 2}
			for _, dest := range loc.Destinations {
				cbData := misc.CallbackData{Action: misc.GO_TO, Payload: fmt.Sprint(dest.GetID())}
				kb.AddButton(dest.Name, cbData.JSON())
			}

			textEdit := tgbotapi.NewEditMessageText(player_C.ChatID, update.CallbackQuery.Message.MessageID, "Go to ...")
			markupEdit := tgbotapi.NewEditMessageReplyMarkup(player_C.ChatID, update.CallbackQuery.Message.MessageID, kb.Generate())
			g.Bot.Send(textEdit)
			g.Bot.Send(markupEdit)
		case "do":
			g.Expeditions.Add(activities.NewExpedition(g.Bot, chrctr))
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
			replyMarkup := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("⚡Start adventure"), tgbotapi.NewKeyboardButton("😐Character")))
			msg.ReplyMarkup = replyMarkup
			g.Bot.Send(msg)
		case misc.SELECT_LOOT_ITEM, misc.PICK_UP_ITEM, misc.PICK_UP_ALL_ITEMS, misc.DISMISS_LOOT:
			g.LootManager.HandleInput(callbackData, update.CallbackQuery.From.ID)
		case misc.EXPEDITION_CONTINUE, misc.EXPEDITION_LEAVE, misc.ACTIVITY_SELECTED, misc.GATHERING_CORRECT, misc.GATHERING_INCORRECT:
			g.Expeditions.HandleInput(update)
		}
	}
}
