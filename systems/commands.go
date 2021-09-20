package systems

import (
	"log"

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

func (b *CommandsSystem) Update(dt float64) {
	for {
		u := bot.GetInstance().PopUpdate("/menu", "/reg")
		if u == nil {
			break
		}
		chrctr := b.characterManager.GetCharacter(u.Message.From.ID)
		if chrctr == nil {
			//TODO Handle missing character
			continue
		}
		switch u.Message.Text {
		case "/menu":
			SendMainMenuKeyboard(chrctr)
		case "/reg":
			log.Println(u.Message.From.ID, u.Message.Chat.ID)
			//.characterManager.CreateCharacter(u.Message.From.ID, u.Message.Chat.ID, "Rustam")
		}
	}
}
