package systems

import (
	"log"
	"strings"

	"github.com/alexglazkov9/survgram/bot"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/interfaces"
)

type CommandsSystem struct {
	currentDt float64

	characterManager interfaces.CharacterHelper
	manager          *entity.Manager
}

func NewCommandsSystem(manager *entity.Manager, characterMgr interfaces.CharacterHelper) *CommandsSystem {
	sys := &CommandsSystem{manager: manager, characterManager: characterMgr}

	return sys
}

func (cs *CommandsSystem) Update(dt float64) {
	for {
		u := bot.GetInstance().PopUpdate("/menu", "/reg")
		if u == nil {
			break
		}
		chrctr := cs.characterManager.GetCharacter(u.Message.From.ID)
		cmd := strings.Split(u.Message.Text, " ")[0]

		switch cmd {
		case "/menu":
			if chrctr != nil {
				SendMainMenuKeyboard(chrctr)
			}
		case "/reg":

			log.Println(u.Message.From.ID, u.Message.Chat.ID)
			if chrctr == nil {
				name := strings.Split(u.Message.Text, " ")[1]
				cs.characterManager.CreateCharacter(u.Message.From.ID, u.Message.Chat.ID, name)
			}
		}
	}
}
