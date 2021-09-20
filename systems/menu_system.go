package systems

import (
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/interfaces"
)

type MenuSystem struct {
	manager         *entity.Manager
	characterHelper interfaces.CharacterHelper
}

func NewMenuSystem(manager *entity.Manager, characterHelper interfaces.CharacterHelper) *MenuSystem {
	sys := &MenuSystem{manager: manager, characterHelper: characterHelper}

	return sys
}

func (ms *MenuSystem) Update(dt float64) {
	// for {
	// 	u := bot.GetInstance().PopUpdate(
	// 		bot.MENU_OPTION_CHARACTER,
	// 		bot.MENU_OPTION_EXPEDITION,
	// 		bot.MENU_OPTION_NPCS,
	// 		bot.MENU_OPTION_HIDEOUT,
	// 	)
	// 	if u == nil {
	// 		break
	// 	}
	// 	chrctr := ms.characterHelper.GetCharacter(u.Message.From.ID)
	// 	if chrctr == nil {
	// 		//TODO Handle missing character
	// 		continue
	// 	}
	// 	menu_C, menu_C_ok := chrctr.GetComponent("MenuComponent").(*components.MenuComponent)
	// 	if !menu_C_ok {
	// 		menu_C.Menus.Push()
	// 	}

	// }

	// for _, e := range ms.manager.QueryEntities("MenuComponent") {

	// }
}
