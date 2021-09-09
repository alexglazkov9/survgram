package activities

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

/* Locations - provides access to all locations */
type Locations struct {
	// map[id of the location]Location - using location id from the file
	// will allow to add and remove locations to/from the JSON file in any order,
	// including skipping ids
	// Any better solution?
	locations map[int]*Location
}

func NewLocations() *Locations {
	var instance = &Locations{}
	instance.locations = fetchAllLocations()
	log.Printf("Locations fetched from the database. Count: %d\n", len(instance.locations))

	return instance
}

//Locations - Returns all locations
func (m Locations) Locations() []Location {
	locs := make([]Location, 0, len(m.locations))
	for _, loc := range m.locations {
		locs = append(locs, *loc)
	}
	return locs
}

//GetLocation - Fetches a start location
func (m Locations) GetStartLocation() Location {
	return m.GetLocation(1)
}

//GetLocation - Fetches a location by its ID
func (m Locations) GetLocation(id int) Location {
	var lctn Location
	for _, loc := range m.locations {
		if id == loc.ID {
			return *loc
		}
	}
	return lctn
}

func fetchAllLocations() map[int]*Location {
	file, _ := ioutil.ReadFile("./resources/locations.json")
	var locs interface{}
	_ = json.Unmarshal([]byte(file), &locs)

	locations := make(map[int]*Location, len(locs.([]interface{})))
	/* Fetch all locations and create a map */
	for _, loc := range locs.([]interface{}) {
		loc_id := int(loc.(map[string]interface{})["id"].(float64))
		loc_name := loc.(map[string]interface{})["name"].(string)
		locations[loc_id] = NewLocation(
			loc_id,
			loc_name,
		)
	}

	/* Assign all destinations */
	//Loop over each location to assign destinations
	for _, loc := range locs.([]interface{}) {
		loc_id := int(loc.(map[string]interface{})["id"].(float64))
		destinations := loc.(map[string]interface{})["destinations"].([]interface{})

		locations[loc_id].Destinations = make([]*Location, len(destinations))
		//Loop over destination ids
		for _, dest := range destinations {
			dest_id := int(dest.(float64))
			destination := locations[dest_id]
			locations[loc_id].Destinations = append(locations[loc_id].Destinations, destination)
		}
	}

	return locations
}
