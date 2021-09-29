package systems

import (
	"github.com/alexglazkov9/survgram/bot"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/interfaces"
)

/* Handles menu at the bottom of the screen built with ReplyKeyboard */
type MenuSystem struct {
	manager         *entity.Manager
	characterHelper interfaces.CharacterHelper
}

func NewMenuSystem(manager *entity.Manager, characterHelper interfaces.CharacterHelper) *MenuSystem {
	sys := &MenuSystem{manager: manager, characterHelper: characterHelper}

	return sys
}

func (ms *MenuSystem) Update(dt float64) {
	for {
		u := bot.GetInstance().PopUpdate(
			bot.MENU_OPTION_CHARACTER,
			bot.MENU_OPTION_EXPEDITION,
			bot.MENU_OPTION_NPCS,
			bot.MENU_OPTION_HIDEOUT,
			bot.MENU_OPTION_MAP,

			bot.CHARACTER_CHARACTER,
			bot.CHARACTER_INVENTORY,
			bot.CHARACTER_SKILLS,

			bot.MENU_BACK,

			bot.EXPEDITION_START,
		)
		if u == nil {
			break
		}
		chrctr := ms.characterHelper.GetCharacter(u.Message.From.ID)
		if chrctr == nil {
			//TODO Handle missing character
			continue
		}

		var menu_C *components.MenuComponent
		if chrctr.HasComponent("MenuComponent") {
			menu_C = chrctr.GetComponent("MenuComponent").(*components.MenuComponent)
		} else {
			menu_C = &components.MenuComponent{}
			menu_C.Menus.Push(GetMainMenu(chrctr).(components.Menu))
			chrctr.AddComponent(menu_C)
		}

		var menu components.Menu
		//Handle going back from any menu level
		if u.Message.Text == bot.MENU_BACK {
			var ok bool
			if len(menu_C.Menus) > 1 {
				menu_C.Menus.Pop()
			}
			menu, ok = menu_C.Menus.Top()
			if !ok {
				continue
			}
			bot.GetInstance().GetBot().Send(menu.Msg)
			continue
		}

		//Handle the selected menu option
		menu, ok := menu_C.Menus.Top()
		if ok {
			action, ok := menu.MenuOptions[u.Message.Text]
			if ok {
				mn := action(chrctr)
				//Functions that do not return menu - return nil
				if mn != nil {
					bot.GetInstance().GetBot().Send(mn.(components.Menu).Msg)
					menu_C.Menus.Push(mn.(components.Menu))
				}
			} else { //resets to main menu if the sent command is not in the list of options
				mn := GetMainMenu(chrctr)
				bot.GetInstance().GetBot().Send(mn.(components.Menu).Msg)
				menu_C.Menus.Clear()
				menu_C.Menus.Push(mn.(components.Menu))
			}
		} else { //resets to main menu if there are no menus on stack
			mn := GetMainMenu(chrctr)
			bot.GetInstance().GetBot().Send(mn.(components.Menu).Msg)
			menu_C.Menus.Clear()
			menu_C.Menus.Push(mn.(components.Menu))
		}
	}
}
