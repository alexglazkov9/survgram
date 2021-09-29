package components

type StatsComponent struct {
	BaseComponent `bson:"-" json:"-"`

	Int int
	Agi int
	Str int
}

func (comp *StatsComponent) GetAgi() int {
	return comp.Agi
}

func (comp *StatsComponent) GetStr() int {
	return comp.Str
}

func (comp *StatsComponent) GetInt() int {
	return comp.Int
}
