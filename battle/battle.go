package battle

import (
	"log"

	"github.com/alexglazkov9/survgram/entity"
)

type Battle struct {
	sideA []*entity.Entity
	sideB []*entity.Entity
}

func (b Battle) Update(dt float64) {
	for i, e := range b.getAllEntities() {
		(*e).Update(dt)
		log.Printf("%d, %d", (*e).GetBasicAttackDamage(), i)
	}
	for _, e := range b.sideA {
		if (*e).IsReadyToAttack() {
			(*e).Attack(b.sideB)
		}
	}
	for _, e := range b.sideB {
		if (*e).IsReadyToAttack() {
			(*e).Attack(b.sideA)
		}
	}
}

func (b Battle) getAllEntities() []*entity.Entity {
	return append(b.sideA, b.sideB...)
}

func (b *Battle) AddToSideA(e entity.Entity) {
	b.sideA = append(b.sideA, &e)
}

func (b *Battle) AddToSideB(e entity.Entity) {
	b.sideB = append(b.sideB, &e)
}
