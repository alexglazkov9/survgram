package game

import (
	"github.com/alexglazkov9/survgram/activities"
	"github.com/alexglazkov9/survgram/database"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/enemies"
	"github.com/alexglazkov9/survgram/items"
	"github.com/alexglazkov9/survgram/items/loot"
	"github.com/alexglazkov9/survgram/systems"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	gl "github.com/kutase/go-gameloop"
)

type Game struct {
	CharacterManager *database.CharacterManager
	Expeditions      *activities.Expeditions
	LootManager      *loot.LootDispenser
	Bot              *tgbotapi.BotAPI
	Engine           *gl.GameLoop
	manager          *entity.Manager

	//Systems
	battleSystem        *systems.BattleSystem
	commandsSystem      *systems.CommandsSystem
	locationSystem      *systems.LocationSystem
	expeditionSystem    *systems.ExpeditionSystem
	lootDispenserSystem *systems.LootDispenserSystem
	hideoutSystem       *systems.HideoutSystem
	menuSystem          *systems.MenuSystem
}

// GetInstance - returns instance of Character Manager
func New(bot *tgbotapi.BotAPI) *Game {
	instance := &Game{}
	instance.manager = entity.NewManager()
	instance.Bot = bot
	instance.CharacterManager = database.NewCharacterManager(database.GetInstance(), activities.GetLocations().GetStartLocation(), instance.manager)
	instance.LootManager = loot.NewLootManager(*bot)
	instance.Expeditions = &activities.Expeditions{LootManager: instance.LootManager, CharacterManager: instance.CharacterManager}

	//Systems
	instance.battleSystem = systems.NewBattleSystem(instance.manager, instance.CharacterManager)
	instance.commandsSystem = systems.NewCommandsSystem(instance.manager, instance.CharacterManager)
	instance.locationSystem = systems.NewLocationSystem(instance.manager, instance.CharacterManager)
	instance.expeditionSystem = systems.NewExpeditionSystem(instance.manager, instance.CharacterManager)
	instance.lootDispenserSystem = systems.NewLootDispenserSystem(instance.manager, instance.CharacterManager)
	instance.hideoutSystem = systems.NewHideoutSystem(instance.manager, instance.CharacterManager)
	instance.menuSystem = systems.NewMenuSystem(instance.manager, instance.CharacterManager)

	instance.Engine = gl.New(30, func(dt float64) {
		instance.Expeditions.Update(dt)
		instance.LootManager.Update(dt)

		instance.battleSystem.Update(dt)
		instance.commandsSystem.Update(dt)
		instance.locationSystem.Update(dt)
		instance.expeditionSystem.Update(dt)
		instance.lootDispenserSystem.Update(dt)
		instance.hideoutSystem.Update(dt)
		instance.menuSystem.Update(dt)
	})

	instance.Engine.Start()
	enemies.GetInstance().SetManager(instance.manager)
	items.GetItemCollection()
	return instance
}

func (g Game) HandleInput(update tgbotapi.Update) {

}
