package database

import (
	"context"
	"log"
	"time"

	"github.com/alexglazkov9/survgram/location"
	"go.mongodb.org/mongo-driver/bson"
)

//GetAllLocations - Returns array of all locations from database
func (d Database) GetAllLocations() []*location.Location {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := d.client.Database("survgram_dev").Collection("locations").Find(ctx, bson.D{})
	if err != nil {
		log.Fatalln(err)
	}
	defer cur.Close(ctx)

	var locations []*location.Location
	for cur.Next(ctx) {
		var location location.Location
		err := cur.Decode(&location)
		if err != nil {
			log.Fatal(err)
		}
		locations = append(locations, &location)
	}

	return locations
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
