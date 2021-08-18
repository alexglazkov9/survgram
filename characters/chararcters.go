package characters

import (
	"log"

	"sync"

	"github.com/alexglazkov9/survgram/character"
	"github.com/alexglazkov9/survgram/database"
)

var once sync.Once

// Manager - Singleton Character Manager that provides characters
type Manager struct {
	characters []*character.Character
}

var instance *Manager

// GetInstance - returns instance of Character Manager
func GetInstance() *Manager {
	once.Do(func() {
		instance = &Manager{}
		instance.characters = database.GetInstance().GetAllCharacters()
		log.Printf("Characters fetched from the database. Count: %d\n", len(instance.characters))
	})

	return instance
}

//Characters - Fetches all characters
func (m Manager) Characters() []*character.Character {
	return m.characters
}

//GetCharacter - Fetches a charcter by telegramID
func (m Manager) GetCharacter(telegramID int) *character.Character {
	for _, chrctr := range m.characters {
		if telegramID == chrctr.TelegramID {
			return chrctr
		}
	}
	return nil
}

//CreateCharacter - Creates a new charcter and writes it to db
func (m Manager) CreateCharacter(telegramID int, name string) bool {
	chrctr := character.New(telegramID, name)
	if database.GetInstance().AddCharacter(chrctr) {
		m.characters = append(m.characters, chrctr)
		return true
	}
	return false
}
