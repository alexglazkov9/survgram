package characters

import (
	"log"

	"github.com/alexglazkov9/survgram/character"
	"github.com/alexglazkov9/survgram/database"
)

// CharacterManager - Singleton Character CharacterManager that provides characters
type CharacterManager struct {
	characters []*character.Character
}

// GetInstance - returns instance of Character Manager
func New() *CharacterManager {
	var instance = &CharacterManager{}
	instance.characters = database.GetInstance().GetAllCharacters()
	log.Printf("Characters fetched from the database. Count: %d\n", len(instance.characters))

	return instance
}

//Characters - Fetches all characters
func (m CharacterManager) Characters() []*character.Character {
	return m.characters
}

//GetCharacter - Fetches a charcter by telegramID
func (m CharacterManager) GetCharacter(telegramID int) *character.Character {
	for _, chrctr := range m.characters {
		if telegramID == chrctr.TelegramID {
			return chrctr
		}
	}
	return nil
}

//CreateCharacter - Creates a new charcter and writes it to db
func (m *CharacterManager) CreateCharacter(telegramID int, chatID int64, name string) bool {
	startLocation := database.GetInstance().GetStartLocation()
	chrctr := character.New(telegramID, chatID, name, startLocation.ID)
	if database.GetInstance().AddCharacter(chrctr) {
		m.characters = append(m.characters, chrctr)
		return true
	}
	return false
}
