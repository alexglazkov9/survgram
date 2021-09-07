package database

import (
	"context"
	"log"
	"time"

	"github.com/alexglazkov9/survgram/interfaces"
	"github.com/alexglazkov9/survgram/location"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//GetAllLocations - Returns array of all locations from database
func (d Database) GetAllLocations() []*location.LocationModel {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	lookup := mongo.Pipeline{
		{{"$lookup",
			bson.D{
				{
					"from", "locations",
				},
				{
					"localField", "destinations",
				},
				{
					"foreignField", "_id",
				},
				{
					"as", "destinations",
				},
			},
		}},
	}
	cur, err := d.client.Database("survgram_dev").Collection("locations").Aggregate(ctx, lookup)
	if err != nil {
		log.Fatalln(err)
	}
	defer cur.Close(ctx)

	var locations []*location.LocationModel
	for cur.Next(ctx) {
		var location location.LocationModel
		err := cur.Decode(&location)
		if err != nil {
			log.Fatal(err)
		}
		locations = append(locations, &location)
	}

	return locations
}

//GetAllLocations - Returns array of all locations from database
func (d Database) GetStartLocation() interfaces.ILocation {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var location *location.Location
	err := d.client.Database("survgram_dev").Collection("locations").FindOne(ctx, bson.M{"name": "Wimborne"}).Decode(&location)
	if err != nil {
		log.Fatal(err)
	}

	return location
}

//AddLocation - Inserts new location into database, returns `true` on success.
func (d Database) AddLocation(l *location.Location) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Inserting new location: " + l.Name)
	_, err := d.client.Database("survgram_dev").Collection("locations").InsertOne(ctx, l)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Location inserted\n")

	return true
}

//UpdateCharacter - Updates character in database, returns `true` on success.
func (d Database) UpdateLocation(l *location.Location) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Updateing location: " + l.Name)
	_, err := d.client.Database("survgram_dev").Collection("locations").UpdateByID(ctx, l.ID, bson.M{"$set": l})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Location updated\n")

	return true
}
