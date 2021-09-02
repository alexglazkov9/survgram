package ability

import (
	"github.com/alexglazkov9/survgram/interfaces"
)

type Ability struct {
	name        string
	tier        int
	Energy_cost int
	Effects     []interfaces.IEffect
}
