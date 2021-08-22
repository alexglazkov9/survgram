package database

import (
	"context"
	"log"
	"time"

	"github.com/alexglazkov9/survgram/character"
	"go.mongodb.org/mongo-driver/bson"
)

//GetAllCharacters - Returns array of all characters from database
func (d Database) GetAllCharacters() []*character.Character {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := d.client.Database("survgram_dev").Collection("characters").Find(ctx, bson.D{})

	if err != nil {
		log.Fatalln(err)
	}
	defer cur.Close(ctx)

	var characters []*character.Character
	for cur.Next(ctx) {
		var chrctr character.Character
		err := cur.Decode(&chrctr)
		if err != nil {
			log.Fatal(err)
		}
		characters = append(characters, &chrctr)
	}

	return characters
}

//AddCharacter - Inserts new character into database, returns `true` on success.
func (d Database) AddCharacter(c *character.Character) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Inserting new character: " + c.Name)
	_, err := d.client.Database("survgram_dev").Collection("characters").InsertOne(ctx, c)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Character inserted\n")

	return true
}

//UpdateCharacter - Updates character in database, returns `true` on success.
func (d Database) UpdateCharacter(c *character.Character) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Updateing character: " + c.Name)
	_, err := d.client.Database("survgram_dev").Collection("characters").UpdateByID(ctx, c.ID, c)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Character updated\n")

	return true
}
