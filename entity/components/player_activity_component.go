package components

type PlayerActivityComponent struct {
	BaseComponent `bson:"-" json:"-"`

	Type string
}
