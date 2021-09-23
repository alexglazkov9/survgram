package systems

import (
	"fmt"
	"strings"

	"github.com/alexglazkov9/survgram/activities"
	"github.com/alexglazkov9/survgram/bot"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/misc"
	"github.com/alexglazkov9/survgram/resources"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func SendMainMenuKeyboard(e *entity.Entity) {
	menu_C := e.GetComponent("MenuComponent").(*components.MenuComponent)
	menu_C.Menus.Clear()
	menu := GetMainMenu(e).(components.Menu)
	menu_C.Menus.Push(menu)
	bot.GetInstance().GetBot().Send(menu.Msg)
}

func GetMainMenu(e *entity.Entity) interface{} {
	player_C := e.GetComponent("PlayerComponent").(*components.PlayerComponent)
	loc := activities.GetLocations().GetLocation(player_C.CurrentLocation)
	txt := fmt.Sprintf("You are at the %s", loc.Name)
	kb := tgbotapi.NewReplyKeyboard()
	kb.Keyboard = append(kb.Keyboard, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(bot.MENU_OPTION_EXPEDITION),
		tgbotapi.NewKeyboardButton(bot.MENU_OPTION_NPCS),
	))
	kb.Keyboard = append(kb.Keyboard, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(bot.MENU_OPTION_CHARACTER),
		tgbotapi.NewKeyboardButton(bot.MENU_OPTION_HIDEOUT),
	))
	kb.Keyboard = append(kb.Keyboard, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(bot.MENU_OPTION_MAP),
	))
	msg := tgbotapi.NewMessage(player_C.ChatID, txt)
	msg.ReplyMarkup = kb
	menu := components.Menu{
		Msg:         msg,
		MenuOptions: make(map[string]func(*entity.Entity) interface{}),
	}
	menu.MenuOptions[bot.MENU_OPTION_EXPEDITION] = GetExpeditionMenu
	//menu.MenuOptions[bot.MENU_OPTION_NPCS] = nil
	menu.MenuOptions[bot.MENU_OPTION_CHARACTER] = GetCharacterMenu
	menu.MenuOptions[bot.MENU_OPTION_HIDEOUT] = GetHideoutMenu
	menu.MenuOptions[bot.MENU_OPTION_MAP] = DisplayMap

	return menu
}

func GetExpeditionMenu(e *entity.Entity) interface{} {
	player_C := e.GetComponent("PlayerComponent").(*components.PlayerComponent)
	loc := activities.GetLocations().GetLocation(player_C.CurrentLocation)
	msgTxt := ""
	if len(loc.PossibleActivities) == 0 { //If no activities, don't show Expedition menu
		msgTxt = "There are no activities available in this location\n"
		msg := tgbotapi.NewMessage(player_C.ChatID, msgTxt)
		bot.GetInstance().GetBot().Send(msg)
		return nil
	}
	msgTxt = "Possible activities in this location:\n"
	for _, act := range loc.PossibleActivities {
		msgTxt += fmt.Sprintf("%.0f%% %s \n", act.ActivityChance*100, strings.Title(string(act.Type)))
	}
	kb := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(bot.EXPEDITION_START),
		tgbotapi.NewKeyboardButton(bot.EXPEDITION_ACTIVITIES_INFO),
		tgbotapi.NewKeyboardButton(bot.MENU_BACK),
	))
	msg := tgbotapi.NewMessage(player_C.ChatID, msgTxt)
	msg.ReplyMarkup = kb

	menu := components.Menu{
		Msg:         msg,
		MenuOptions: make(map[string]func(*entity.Entity) interface{}),
	}
	menu.MenuOptions[bot.EXPEDITION_START] = StartExpedition
	//menu.MenuOptions[bot.EXPEDITION_ACTIVITIES_INFO] = nil
	return menu
}

func GetHideoutMenu(e *entity.Entity) interface{} {
	player_C := e.GetComponent("PlayerComponent").(*components.PlayerComponent)
	msg := tgbotapi.NewPhotoUpload(player_C.ChatID, "./resources/images/hideout.png")
	msg.Caption = resources.HIDEOUT_WELCOME_TEXT
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(bot.HIDEOUT_STORAGE),
			tgbotapi.NewKeyboardButton(bot.HIDEOUT_CRAFTING_TABLE),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(bot.MENU_BACK),
		),
	)
	msg.ParseMode = "markdown"
	msg.ReplyMarkup = kb
	menu := components.Menu{
		Msg:         msg,
		MenuOptions: make(map[string]func(*entity.Entity) interface{}),
	}
	menu.MenuOptions[bot.HIDEOUT_STORAGE] = nil
	menu.MenuOptions[bot.HIDEOUT_CRAFTING_TABLE] = nil
	return menu
}

func GetCharacterMenu(e *entity.Entity) interface{} {
	player_C := e.GetComponent("PlayerComponent").(*components.PlayerComponent)
	msg := tgbotapi.NewMessage(player_C.ChatID, "View and customize your character here!")
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(bot.CHARACTER_CHARACTER),
			tgbotapi.NewKeyboardButton(bot.CHARACTER_INVENTORY),
			tgbotapi.NewKeyboardButton(bot.CHARACTER_SKILLS),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(bot.MENU_BACK),
		),
	)
	msg.ParseMode = "markdown"
	msg.ReplyMarkup = kb
	menu := components.Menu{
		Msg:         msg,
		MenuOptions: make(map[string]func(*entity.Entity) interface{}),
	}
	//menu.MenuOptions[bot.CHARACTER_CHARACTER] = nil
	menu.MenuOptions[bot.CHARACTER_INVENTORY] = DisplayInventory
	//menu.MenuOptions[bot.CHARACTER_SKILLS] = nil
	return menu
}

func GetExpeditionQuickMenu(e *entity.Entity) interface{} {
	player_C := e.GetComponent("PlayerComponent").(*components.PlayerComponent)
	loc := activities.GetLocations().GetLocation(player_C.CurrentLocation)
	msg := tgbotapi.NewMessage(player_C.ChatID, fmt.Sprintf("%s at %s", resources.EXPEDITION_ENTERED_TEXT, loc.Name))
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(bot.EXPEDITION_RUNNING_MENU_CHARACTER),
			tgbotapi.NewKeyboardButton(bot.EXPEDITION_RUNNING_MENU_INVENTORY),
			tgbotapi.NewKeyboardButton(bot.EXPEDITION_RUNNING_MENU_ABILITIES),
		),
	)
	msg.ReplyMarkup = kb
	menu := components.Menu{
		Msg:         msg,
		MenuOptions: make(map[string]func(*entity.Entity) interface{}),
	}
	//menu.MenuOptions[bot.EXPEDITION_RUNNING_MENU_CHARACTER] = nil
	menu.MenuOptions[bot.EXPEDITION_RUNNING_MENU_INVENTORY] = DisplayInventory
	//menu.MenuOptions[bot.EXPEDITION_RUNNING_MENU_ABILITIES] = nil
	return menu
}

func GetQuickMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(bot.EXPEDITION_RUNNING_MENU_CHARACTER),
			tgbotapi.NewKeyboardButton(bot.EXPEDITION_RUNNING_MENU_INVENTORY),
			tgbotapi.NewKeyboardButton(bot.EXPEDITION_RUNNING_MENU_ABILITIES),
		),
	)
	return kb
}

func StartExpedition(e *entity.Entity) interface{} {
	player_C := e.GetComponent("PlayerComponent").(*components.PlayerComponent)
	loc := activities.GetLocations().GetLocation(player_C.CurrentLocation)
	if e.HasComponent("PlayerActivityComponent") { //check if player is busy
		return nil
	}

	//Start expedition
	expdtnComp := &components.ExpeditionComponent{
		State:    components.STARTING,
		Players:  []*entity.Entity{},
		IsReady:  true,
		Location: loc,
		Messages: make(map[int]tgbotapi.Message),
	}
	expdtn := e.Manager.NewEntity()
	expdtnComp.AddPlayer(e)
	expdtn.AddComponent(expdtnComp)
	return GetExpeditionQuickMenu(e)
}

func DisplayMap(e *entity.Entity) interface{} {
	player_C, _ := e.GetComponent("PlayerComponent").(*components.PlayerComponent)

	//Add destinations to the keyboard
	loc := activities.GetLocations().GetLocation(player_C.CurrentLocation)
	kb := misc.TGInlineKeyboard{Columns: 2, IsClosable: true}
	for _, dest := range loc.Destinations {
		cbData := misc.CallbackData{Action: misc.GO_TO, Payload: fmt.Sprint(dest.GetID())}
		kb.AddButton(dest.Name, cbData.JSON())
	}
	msg := tgbotapi.NewMessage(
		player_C.ChatID,
		"This is a map of Survgram. Nice image of the map with player's current position highlighted.\n\nClick a destination below to go there.",
	)
	msg.ReplyMarkup = kb.Generate()
	bot.GetInstance().GetBot().Send(msg)
	return nil
}

func DisplayInventory(e *entity.Entity) interface{} {
	inventoryWindow_C := &components.InventoryWindowComponent{
		IsSent:         false,
		SelectedItemID: nil,
	}
	e.AddComponent(inventoryWindow_C)
	return nil
}
