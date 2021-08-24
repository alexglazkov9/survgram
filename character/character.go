package character

import (
	"math/rand"

	"github.com/alexglazkov9/survgram/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const ATTACK_SPEED float64 = 1.5

// Character sturct
type Character struct {
	ID              primitive.ObjectID `bson:"_id"`
	TelegramID      int
	ChatID          int64
	Name            string
	HealthPoints    int
	CurrentLocation primitive.ObjectID
	AttackDamage    int

	nextAttackTime float64
}

// New - Creates a new character
func New(telegramID int, chatID int64, name string, current_location primitive.ObjectID) *Character {
	return &Character{ID: primitive.NewObjectID(), TelegramID: telegramID, ChatID: chatID, Name: name, HealthPoints: 10, CurrentLocation: current_location, AttackDamage: 7}
}

func (c *Character) Update(dt float64) {
	c.nextAttackTime += dt
}

func (c Character) GetBasicAttackDamage() int {
	return c.AttackDamage
}

func (c *Character) TakeDamage(dmg int) int {
	c.HealthPoints -= dmg
	return dmg
}

func (c Character) IsReadyToAttack() bool {
	return c.nextAttackTime >= ATTACK_SPEED
}

func (c Character) Attack(targets []*entity.Entity) {
	target := targets[rand.Intn(len(targets))]

	(*target).TakeDamage(c.GetBasicAttackDamage())
}
