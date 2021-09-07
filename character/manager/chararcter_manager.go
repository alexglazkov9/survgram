package charactermanager

import (
	"log"

	"github.com/alexglazkov9/survgram/components"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/interfaces"
)

// CharacterManager - Singleton Character CharacterManager that provides characters
type CharacterManager struct {
	characters []*entity.Entity

	database interfaces.IDatabase
}

// GetInstance - returns instance of Character Manager
func New(database interfaces.IDatabase) *CharacterManager {
	var instance = &CharacterManager{database: database}
	instance.characters = instance.database.GetAllCharacters()
	log.Printf("Characters fetched from the database. Count: %d\n", len(instance.characters))

	return instance
}

//Characters - Fetches all characters
func (m CharacterManager) Characters() []*entity.Entity {
	return m.characters
}

//GetCharacter - Fetches a charcter by telegramID
func (m CharacterManager) GetCharacter(telegramID int) *entity.Entity {
	for _, chrctr := range m.characters {
		player_C := chrctr.GetComponent("PlayerComponent").(*components.PlayerComponent)
		if telegramID == player_C.TelegramID {
			return chrctr
		}
	}
	return nil
}

//CreateCharacter - Creates a new charcter and writes it to db
func (m *CharacterManager) CreateCharacter(telegramID int, chatID int64, name string) bool {

	chrctr := m.newCharacter(telegramID, chatID, name)
	if m.database.AddCharacter(chrctr) {
		chrctr.AddComponent(&components.EffectsComponent{})
		m.characters = append(m.characters, chrctr)
		return true
	}
	return false
}

func (m *CharacterManager) UpdateCharacter(e *entity.Entity) {
	m.database.UpdateCharacter(e)
}

func (m *CharacterManager) newCharacter(telegramID int, chatID int64, name string) *entity.Entity {
	startLocation := m.database.GetStartLocation()

	chrctr := entity.New()
	chrctr.AddComponent(&components.PlayerComponent{
		TelegramID:      telegramID,
		ChatID:          chatID,
		CurrentLocation: startLocation.GetID(),
	})
	chrctr.AddComponent(&components.NameComponent{
		Name: name,
	})
	chrctr.AddComponent(&components.AttackComponent{
		AttackDamage: 5,
	})
	chrctr.AddComponent(&components.HealthComponent{
		MaxHealthPoints: 100,
		HealthPoints:    100,
	})
	chrctr.AddComponent(&components.AbilityComponent{})
	chrctr.AddComponent(&components.InventoryComponent{})

	return chrctr
}
