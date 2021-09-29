package activities

import (
	"math/rand"
)

type Location struct {
	ID                 int
	Name               string
	Destinations       []*Location
	PossibleActivities []ActivityConfig
}

/* Creates activity from one of the ActivityConfig assigned to this location and returns it */
func (l Location) GetActivity() ActivityConfig {
	act_index := rand.Intn(len(l.PossibleActivities))
	return l.PossibleActivities[act_index]
}

func (l Location) GetID() int {
	return l.ID
}

// NewLocation - Creates a new character
func NewLocation(id int, name string) *Location {
	return &Location{ID: id, Name: name, Destinations: make([]*Location, 0)}
}
