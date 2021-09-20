package interfaces

import "github.com/alexglazkov9/survgram/entity"

type IDatabase interface {
	UpdateCharacter(c *entity.Entity) bool
	GetAllCharacters(mgr *entity.Manager) []*entity.Entity
	AddCharacter(c *entity.Entity) bool
}

type CharacterHelper interface {
	GetCharacter(telegramID int) *entity.Entity
	CreateCharacter(telegramID int, chatID int64, name string) bool
}
