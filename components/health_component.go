package components

import (
	"log"

	"github.com/alexglazkov9/survgram/entity"
)

type HealthComponent struct {
	Parent *entity.Entity `bson:"-"`

	MaxHealthPoints int
	HealthPoints    int
}

func (hc HealthComponent) Update(dt float64) {

}

func (hc *HealthComponent) ApplyDamage(dmg int) int {
	hc.HealthPoints -= dmg
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
