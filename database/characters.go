package database

import (
	"context"
	"log"
	"time"

	"github.com/alexglazkov9/survgram/ability"
	"github.com/alexglazkov9/survgram/components"
	"github.com/alexglazkov9/survgram/effect"
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/interfaces"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//GetAllCharacters - Returns array of all characters from database
func (d Database) GetAllCharacters() []*entity.Entity {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := d.client.Database("survgram_dev").Collection("characters").Find(ctx, bson.D{})

	if err != nil {
		log.Fatalln(err)
	}
	defer cur.Close(ctx)

	var characters []*entity.Entity
	for cur.Next(ctx) {
		chrctr := entity.New()
		var rawChrctr bson.M
		err := cur.Decode(&rawChrctr)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(rawChrctr)

		//Manually map fields to the character
		chrctr.ID = rawChrctr["_id"].(primitive.ObjectID)
		for key, element := range rawChrctr["components"].(primitive.M) {
			switch key {
			case "AttackComponent":
				chrctr.AddComponent(&components.AttackComponent{
					AttackDamage: int(element.(primitive.M)["attackdamage"].(int32)),
				})
			case "HealthComponent":
				chrctr.AddComponent(&components.HealthComponent{
					MaxHealthPoints: int(element.(primitive.M)["maxhealthpoints"].(int32)),
					HealthPoints:    int(element.(primitive.M)["healthpoints"].(int32)),
				})
			case "PlayerComponent":
				chrctr.AddComponent(&components.PlayerComponent{
					TelegramID:      int(element.(primitive.M)["telegramid"].(int32)),
					ChatID:          element.(primitive.M)["chatid"].(int64),
					CurrentLocation: element.(primitive.M)["currentlocation"].(primitive.ObjectID),
				})
			case "NameComponent":
				chrctr.AddComponent(&components.NameComponent{
					Name: element.(primitive.M)["name"].(string),
				})
			case "AbilityComponent":
				if element.(primitive.M)["ability"] != nil {
					chrctr.AddComponent(&components.AbilityComponent{
						Ability: ability.Ability{
							Energy_cost: int(element.(primitive.M)["ability"].(primitive.M)["energy_cost"].(int32)),
							Effects:     []interfaces.IEffect{},
						},
					})
				} else {
					chrctr.AddComponent(&components.AbilityComponent{})
				}

			}
		}
		chrctr.AddComponent(&components.EffectsComponent{})

		fireball := &ability.Ability{
			Energy_cost: 15,
			Effects: []interfaces.IEffect{
				&effect.MagicalDamageEffect{
					Source: chrctr,
					Damage: 15,
				},
				&effect.PeriodicEffect{
					Source:   chrctr,
					Damage:   3,
					Duration: 15,
					Period:   3,
				},
			},
		}
		chrctr.GetComponent("AbilityComponent").(*components.AbilityComponent).Ability = *fireball
		characters = append(characters, chrctr)
	}

	return characters
}

//AddCharacter - Inserts new character into database, returns `true` on success.
func (d Database) AddCharacter(c *entity.Entity) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	name_C := c.GetComponent("NameComponent").(*components.NameComponent)
	log.Println("Inserting new character: " + name_C.Name)
	_, err := d.client.Database("survgram_dev").Collection("characters").InsertOne(ctx, *c)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Character inserted\n")

	return true
}

//UpdateCharacter - Updates character in database, returns `true` on success.
func (d Database) UpdateCharacter(c *entity.Entity) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	name_C := c.GetComponent("NameComponent").(*components.NameComponent)
	log.Println("Updateing character: " + name_C.Name)
	_, err := d.client.Database("survgram_dev").Collection("characters").UpdateByID(ctx, c.ID, c)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Character updated\n")

	return true
}
