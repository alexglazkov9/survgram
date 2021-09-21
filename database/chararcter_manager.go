package database

import (
	"log"

	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/interfaces"
	"github.com/alexglazkov9/survgram/items"
)

// CharacterManager - CharacterManager that provides characters
type CharacterManager struct {
	startLocation interfaces.ILocation
	manager       *entity.Manager

	database interfaces.IDatabase
}

// GetInstance - returns instance of Character Manager
func NewCharacterManager(database interfaces.IDatabase, startLocation interfaces.ILocation, manager *entity.Manager) *CharacterManager {
	var instance = &CharacterManager{database: database, manager: manager}
	instance.startLocation = startLocation

	characters := instance.database.GetAllCharacters(manager)
	for _, c := range characters {
		c.SetManager(instance.manager)
	}

	log.Printf("Characters fetched from the database. Count: %d\n", len(characters))

	return instance
}

//Characters - Fetches all characters
func (m CharacterManager) Characters() []*entity.Entity {
	return m.manager.QueryEntities("PlayerComponent")
}

//GetCharacter - Fetches a charcter by telegramID
func (m CharacterManager) GetCharacter(telegramID int) *entity.Entity {
	players := m.manager.QueryEntities("PlayerComponent")
	for _, p := range players {
		player_C := p.GetComponent("PlayerComponent").(*components.PlayerComponent)
		if player_C.TelegramID == telegramID {
			return p
		}
	}
	return nil
}

//CreateCharacter - Creates a new charcter and writes it to db
func (m *CharacterManager) CreateCharacter(telegramID int, chatID int64, name string) bool {

	chrctr := m.newCharacter(telegramID, chatID, name)
	if m.database.AddCharacter(chrctr) {
		chrctr.AddComponent(&components.EffectsComponent{})
		return true
	}
	return false
}

func (m *CharacterManager) UpdateCharacter(e *entity.Entity) {
	m.database.UpdateCharacter(e)
}

func (m *CharacterManager) newCharacter(telegramID int, chatID int64, name string) *entity.Entity {
	startLocation := m.startLocation

	chrctr := m.manager.NewEntity()
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
	chrctr.AddComponent(&components.InventoryComponent{
		Items: make([]*items.ItemBundle, 0),
		Slots: 6,
	})

	return chrctr
}
