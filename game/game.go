package game

import (
	"github.com/alexglazkov9/survgram/activities"
	"github.com/alexglazkov9/survgram/database"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/enemies"
	"github.com/alexglazkov9/survgram/items"
	"github.com/alexglazkov9/survgram/systems"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	gl "github.com/kutase/go-gameloop"
)

type Game struct {
	CharacterManager *database.CharacterManager
	Bot              *tgbotapi.BotAPI
	Engine           *gl.GameLoop
	manager          *entity.Manager

	//Systems
	battleSystem        *systems.BattleSystem
	gatheringSystem     *systems.GatheringSystem
	commandsSystem      *systems.CommandsSystem
	locationSystem      *systems.LocationSystem
	expeditionSystem    *systems.ExpeditionSystem
	lootDispenserSystem *systems.LootDispenserSystem
	hideoutSystem       *systems.HideoutSystem
	menuSystem          *systems.MenuSystem
	windowsSystem       *systems.WindowsSystem
}

// GetInstance - returns instance of Character Manager
func New(bot *tgbotapi.BotAPI) *Game {
	instance := &Game{}
	instance.manager = entity.NewManager()
	instance.Bot = bot
	instance.CharacterManager = database.NewCharacterManager(
		database.GetInstance(),
		activities.GetLocations().GetStartLocation(),
		instance.manager,
	)

	//Systems
	instance.battleSystem = systems.NewBattleSystem(instance.manager, instance.CharacterManager)
	instance.gatheringSystem = systems.NewGatheringSystem(instance.manager, instance.CharacterManager)
	instance.commandsSystem = systems.NewCommandsSystem(instance.manager, instance.CharacterManager)
	instance.locationSystem = systems.NewLocationSystem(instance.manager, instance.CharacterManager)
	instance.expeditionSystem = systems.NewExpeditionSystem(instance.manager, instance.CharacterManager)
	instance.lootDispenserSystem = systems.NewLootDispenserSystem(instance.manager, instance.CharacterManager)
	instance.hideoutSystem = systems.NewHideoutSystem(instance.manager, instance.CharacterManager)
	instance.menuSystem = systems.NewMenuSystem(instance.manager, instance.CharacterManager)
	instance.windowsSystem = systems.NewWindowsSystem(instance.manager, instance.CharacterManager)

	instance.Engine = gl.New(30, func(dt float64) {
		//Systems
		instance.battleSystem.Update(dt)
		instance.gatheringSystem.Update(dt)
		instance.commandsSystem.Update(dt)
		instance.locationSystem.Update(dt)
		instance.expeditionSystem.Update(dt)
		instance.lootDispenserSystem.Update(dt)
		instance.hideoutSystem.Update(dt)
		instance.menuSystem.Update(dt)
		instance.windowsSystem.Update(dt)
	})

	enemies.GetInstance().SetManager(instance.manager)
	items.GetItemCollection()

	instance.Engine.Start()

	return instance
}
