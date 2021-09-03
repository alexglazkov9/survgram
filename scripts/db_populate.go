package main

import (
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	//var locations []*location.Location

	// location1 := location.New("Test1")
	// location2 := location.New("Test2")
	// location3 := location.New("Test3")
	// location1.Destinations = append(location1.Destinations, location2.ID)
	// location1.Destinations = append(location1.Destinations, location3.ID)
	// location2.Destinations = append(location2.Destinations, location3.ID)
	// database.GetInstance().AddLocation(location1)
	// database.GetInstance().AddLocation(location2)
	// database.GetInstance().AddLocation(location3)

	// for _, loc := range database.GetInstance().GetAllLocations() {
	// 	log.Println(loc)
	// 	loc.Name += "1"
	// 	database.GetInstance().UpdateLocation(loc)
	// }

	//Generate Map
	// wimborne := location.New("Wimborne")
	// thecollcavern := location.New("The Collapsing Cavern")
	// whispvalley := location.New("Whispering Valley")
	// oldforest := location.New("Old Forest")
	// buxpier := location.New("Buxbourg's Pier")
	// buxlake := location.New("Buxbourg Lake")
	// neverharbor := location.New("Neverharbor")
	// stokefallpier := location.New("Stokefall's Pier")
	// stokefallfalls := location.New("Stokefall Falls")

	// wimborne.Destinations = append(wimborne.Destinations, thecollcavern.ID, oldforest.ID, buxpier.ID, stokefallpier.ID, thecollcavern.ID)
	// thecollcavern.Destinations = append(thecollcavern.Destinations, wimborne.ID)
	// whispvalley.Destinations = append(whispvalley.Destinations, wimborne.ID)
	// oldforest.Destinations = append(oldforest.Destinations, buxpier.ID, stokefallpier.ID, wimborne.ID)
	// buxpier.Destinations = append(buxpier.Destinations, wimborne.ID, oldforest.ID, buxlake.ID)
	// buxlake.Destinations = append(buxlake.Destinations, buxpier.ID, neverharbor.ID)
	// neverharbor.Destinations = append(neverharbor.Destinations, buxlake.ID)
	// stokefallpier.Destinations = append(stokefallpier.Destinations, oldforest.ID, wimborne.ID, stokefallfalls.ID)

	// database.GetInstance().AddLocation(wimborne)
	// database.GetInstance().AddLocation(thecollcavern)
	// database.GetInstance().AddLocation(whispvalley)
	// database.GetInstance().AddLocation(oldforest)
	// database.GetInstance().AddLocation(buxpier)
	// database.GetInstance().AddLocation(buxlake)
	// database.GetInstance().AddLocation(neverharbor)
	// database.GetInstance().AddLocation(stokefallpier)
	// database.GetInstance().AddLocation(stokefallfalls)

}
