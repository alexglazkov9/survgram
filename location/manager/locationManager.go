package manager

import (
	"log"

	"github.com/alexglazkov9/survgram/database"
	"github.com/alexglazkov9/survgram/location"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LocationManager - provides locations
type LocationManager struct {
	locations []*location.LocationModel
}

// GetInstance - returns instance of LocationManager
func New() *LocationManager {
	var instance = &LocationManager{}
	instance.locations = database.GetInstance().GetAllLocations()
	log.Printf("Locations fetched from the database. Count: %d\n", len(instance.locations))

	return instance
}

//Locations - Fetches all locations
func (m LocationManager) Locations() []*location.LocationModel {
	return m.locations
}

//GetLocation - Fetches a location by its ID
func (m LocationManager) GetLocation(objectID primitive.ObjectID) *location.LocationModel {
	for _, loc := range m.locations {
		if objectID == loc.ID {
			return loc
		}
	}
	return nil
}
