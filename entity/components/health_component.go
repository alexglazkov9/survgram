package components

import (
	"log"

	"github.com/alexglazkov9/survgram/entity"
)

type HealthComponent struct {
	BaseComponent `bson:"-" json:"-"`

	HealthPoints    int
	MaxHealthPoints int
}

func (hc HealthComponent) Update(dt float64) {

}

func (hc *HealthComponent) ApplyDamage(dmg int) int {
	hc.HealthPoints -= dmg
	if hc.HealthPoints < 0 {
		hc.HealthPoints = 0
	}
	log.Printf("HP: %d/%d", hc.HealthPoints, hc.MaxHealthPoints)
	return dmg
}

func (hc HealthComponent) GetHealth() int {
	return hc.HealthPoints
}

func (hc HealthComponent) GetMaxHealth() int {
	return hc.MaxHealthPoints
}

func (hc HealthComponent) IsAlive() bool {
	return hc.HealthPoints > 0
}

func (hc *HealthComponent) GetHealthComponent() *HealthComponent {
	return hc
}

func (ac *HealthComponent) Clone() entity.IComponent {
	copy := *ac
	return &copy
}
