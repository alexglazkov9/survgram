package components

import "github.com/alexglazkov9/survgram/entity"

type LogComponent struct {
	BaseComponent `bson:"-" json:"-"`

	Log *[]string
}

func (ac *LogComponent) Update(dt float64) {
}

func (lc *LogComponent) AddRecord(record string) {
	*lc.Log = append(*lc.Log, record)
}

func (ac *LogComponent) Clone() entity.IComponent {
	copy := *ac
	return &copy
}
