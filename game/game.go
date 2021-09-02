package game

import (
	"log"

	"github.com/alexglazkov9/survgram/activity"
	activitymanager "github.com/alexglazkov9/survgram/activity/manager"
	charactermanager "github.com/alexglazkov9/survgram/character/manager"
	"github.com/alexglazkov9/survgram/components"
	locationmanager "github.com/alexglazkov9/survgram/location/manager"
	"github.com/alexglazkov9/survgram/misc"
	"go.mongodb.org/mongo-driver/bson/primitive"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	gameLoop "github.com/kutase/go-gameloop"
	gl "github.com/kutase/go-gameloop"
)

type Game struct {
	LocationManager  *locationmanager.LocationManager
	CharacterManager *charactermanager.CharacterManager
	ActivitiyManager *activitymanager.ActivityManager
	Bot              *tgbotapi.BotAPI
	engine           *gl.GameLoop
}

// New - Creates a new game isntance
func New(bot *tgbotapi.BotAPI) *Game {
	game := &Game{}
	game.Bot = bot
	game.Init()
	return game
}

func (g *Game) Init() {
	g.CharacterManager = charactermanager.New()
	g.LocationManager = locationmanager.New()
	g.ActivitiyManager = &activitymanager.ActivityManager{}
	g.engine = gameLoop.New(30, func(delta float64) {
		//log.Println("tick")
		//log.Println(delta)
		g.ActivitiyManager.Update(delta)
	})
	g.engine.Start()
}

func (g Game) HandleInput(update tgbotapi.Update) {
	callbackData := misc.CallbackData{}
	callbackData.FromJSON(update.CallbackQuery.Data)

	chrctr := g.CharacterManager.GetCharacter(update.CallbackQuery.From.ID)
	switch update.CallbackQuery.Data {
	case "goto":
		chrctr := g.CharacterManager.GetCharacter(update.CallbackQuery.From.ID)
		player_C := chrctr.GetComponent("PlayerComponent").(components.PlayerComponent)

		var buttons tgbotapi.InlineKeyboardMarkup

		loc := g.LocationManager.GetLocation(player_C.CurrentLocation)

		//Add destinations to the keyboard
		var row []tgbotapi.InlineKeyboardButton
		i := 0
		for _, dest := range loc.Destinations {
			cbData := misc.CallbackData{Action: misc.GO_TO, Payload: dest.ID.Hex()}
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
		g.ActivitiyManager.Add(activity.New(*g.Bot, chrctr))
	}

	//Moves character to a new location
	switch callbackData.Action {
	case misc.GO_TO:
		chrctr := g.CharacterManager.GetCharacter(update.CallbackQuery.From.ID)
		player_C := chrctr.GetComponent("PlayerComponent").(components.PlayerComponent)

		objID, err := primitive.ObjectIDFromHex(callbackData.Payload)
		if err != nil {
			log.Fatal(err)
		}
		player_C.CurrentLocation = objID
		msg := tgbotapi.NewMessage(player_C.ChatID, "You have reached the new location")
		g.Bot.Send(msg)
	}
}
