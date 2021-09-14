package effect

import (
	"fmt"

	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
)

type BaseEffect struct {
	Type string
}

type MagicalDamageEffect struct {
	BaseEffect

	Source *entity.Entity `bson:"-"`

	Damage int
}

func (ef *MagicalDamageEffect) Update(dt float64) {
}

func (ef *MagicalDamageEffect) Apply(target *entity.Entity) bool {
	target.GetComponent("HealthComponent").(*components.HealthComponent).ApplyDamage(ef.Damage)

	log_C := target.GetComponent("LogComponent").(*components.LogComponent)
	name_C := target.GetComponent("NameComponent").(*components.NameComponent)
	source_name_C := ef.Source.GetComponent("NameComponent").(*components.NameComponent)
	log_C.AddRecord(fmt.Sprintf("%s casted spell on %s for %d damage", source_name_C.GetName(), name_C.GetName(), ef.Damage))

	return true
}

type PeriodicEffect struct {
	BaseEffect

	Source *entity.Entity `bson:"-"`

	Damage   int
	Duration float64
	Period   float64

	timeTotal  float64
	timePeriod float64
}

func (ef *PeriodicEffect) Update(dt float64) {
	ef.timeTotal += dt
	ef.timePeriod += dt
}

func (ef *PeriodicEffect) Apply(target *entity.Entity) bool {
	if ef.timeTotal >= ef.Duration {
		return true
	}
	if ef.timePeriod >= ef.Period {
		ef.timePeriod -= ef.Period
		target.GetComponent("HealthComponent").(*components.HealthComponent).ApplyDamage(ef.Damage)

		//TODO figure out whether to remove effects upon battle end or wait till they expire
		log_C, ok := target.GetComponent("LogComponent").(*components.LogComponent)
		if ok {
			name_C := target.GetComponent("NameComponent").(*components.NameComponent)
			source_name_C := ef.Source.GetComponent("NameComponent").(*components.NameComponent)
			log_C.AddRecord(fmt.Sprintf("%s burned %s for %d damage", source_name_C.GetName(), name_C.GetName(), ef.Damage))
		}
	}
	return false
}
