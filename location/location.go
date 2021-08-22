package location

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string
	Destinations []primitive.ObjectID
}

type LocationModel struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string
	Destinations []Location
}

// New - Creates a new character
func New(name string) *Location {
	return &Location{ID: primitive.NewObjectID(), Name: name}
}

// // New - Creates a new character
// func (l Location) Save() {
// 	database.GetInstance().UpdateLocation(l)
// }
