package database

import (
	"context"
	"log"
	"time"

	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/combat"
	"github.com/alexglazkov9/survgram/entity/combat/effect"
	"github.com/alexglazkov9/survgram/entity/components"
	"github.com/alexglazkov9/survgram/interfaces"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//GetAllCharacters - Returns array of all characters from database
func (d Database) GetAllCharacters(manager *entity.Manager) []*entity.Entity {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := d.client.Database("survgram_dev").Collection("characters").Find(ctx, bson.D{})

	if err != nil {
		log.Fatalln(err)
	}
	defer cur.Close(ctx)

	var characters []*entity.Entity
	for cur.Next(ctx) {
		var rawChrctr bson.M
		err := cur.Decode(&rawChrctr)
		if err != nil {
			log.Fatal(err)
		}

		chrctr := parseCharacter(rawChrctr, manager)
		characters = append(characters, chrctr)
	}

	return characters
}

//AddCharacter - Inserts new character into database, returns `true` on success.
func (d Database) AddCharacter(c *entity.Entity) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	name_C := c.GetComponent("NameComponent").(*components.NameComponent)
	cmpnts := getComponentsMap(c)

	log.Println("Inserting new character: " + name_C.Name)
	_, err := d.client.Database("survgram_dev").Collection("characters").InsertOne(ctx, bson.D{{"components", cmpnts}})
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

	cmpnts := getComponentsMap(c)

	log.Println("Updating character")
	_, err := d.client.Database("survgram_dev").Collection("characters").UpdateByID(ctx, c.MongoID, bson.D{
		{"$set", bson.D{{"components", cmpnts}}},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Character updated\n")

	return true
}

/* Manually map fields to the character */
func parseCharacter(raw primitive.M, manager *entity.Manager) *entity.Entity {
	chrctr := manager.NewEntity()
	chrctr.MongoID = raw["_id"].(primitive.ObjectID)
	for key, element := range raw["components"].(primitive.M) {
		// To benefit from built-in parser, primitive.M component is converted
		// to bson and later converted into struct of the proper component type
		bsonElement, _ := bson.Marshal(element)

		switch key {
		case "AttackComponent":
			var comp components.AttackComponent
			bson.Unmarshal(bsonElement, &comp)
			chrctr.AddComponent(&comp)
		case "HealthComponent":
			var comp components.HealthComponent
			bson.Unmarshal(bsonElement, &comp)
			chrctr.AddComponent(&comp)
		case "PlayerComponent":
			var comp components.PlayerComponent
			bson.Unmarshal(bsonElement, &comp)
			chrctr.AddComponent(&comp)
		case "NameComponent":
			var comp components.NameComponent
			bson.Unmarshal(bsonElement, &comp)
			chrctr.AddComponent(&comp)
		case "AbilityComponent":
			if element.(primitive.M)["ability"] != nil {
				var comp components.AbilityComponent
				bson.Unmarshal(bsonElement, &comp)
				chrctr.AddComponent(&comp)
			} else {
				chrctr.AddComponent(&components.AbilityComponent{})
			}
		case "InventoryComponent":
			var comp components.InventoryComponent
			bson.Unmarshal(bsonElement, &comp)
			chrctr.AddComponent(&comp)
		case "PlayerLocationComponent":
			var comp components.PlayerLocationComponent
			bson.Unmarshal(bsonElement, &comp)
			chrctr.AddComponent(&comp)
		case "StatsComponent":
			var comp components.StatsComponent
			bson.Unmarshal(bsonElement, &comp)
			chrctr.AddComponent(&comp)
		case "PlayerEquipmentComponent":
			var comp components.PlayerEquipmentComponent
			comp.Equipment = make([]*components.ItemID, 10)
			bson.Unmarshal(bsonElement, &comp)
			chrctr.AddComponent(&comp)
		}
	}

	chrctr.AddComponent(&components.MenuComponent{})
	chrctr.AddComponent(&components.EffectsComponent{})

	/*TODO REmove this*/
	fireball := &combat.Ability{
		Energy_cost: 15,
		Effects: []interfaces.IEffect{
			&effect.MagicalDamageEffect{
				BaseEffect: effect.BaseEffect{
					Type: "MagicalDamageEffect",
				},
				Source: chrctr,
				Damage: 15,
			},
			&effect.PeriodicEffect{
				BaseEffect: effect.BaseEffect{
					Type: "PeriodicEffect",
				},
				Source:   chrctr,
				Damage:   3,
				Duration: 15,
				Period:   3,
			},
		},
	}
	chrctr.GetComponent("AbilityComponent").(*components.AbilityComponent).Ability = fireball
	/*END TODO Remove this*/

	return chrctr
}

func getComponentsMap(e *entity.Entity) map[string]entity.IComponent {
	name_C := e.GetComponent("NameComponent").(*components.NameComponent)
	health_C := e.GetComponent("HealthComponent").(*components.HealthComponent)
	player_C := e.GetComponent("PlayerComponent").(*components.PlayerComponent)
	inventory_C := e.GetComponent("InventoryComponent").(*components.InventoryComponent)
	attack_C := e.GetComponent("AttackComponent").(*components.AttackComponent)
	ability_C := e.GetComponent("AbilityComponent").(*components.AbilityComponent)
	location_C := e.GetComponent("PlayerLocationComponent").(*components.PlayerLocationComponent)
	stats_C := e.GetComponent("StatsComponent").(*components.StatsComponent)
	equipment_C := e.GetComponent("PlayerEquipmentComponent").(*components.PlayerEquipmentComponent)

	var cmpnts map[string]entity.IComponent
	cmpnts = make(map[string]entity.IComponent)
	cmpnts["NameComponent"] = name_C
	cmpnts["HealthComponent"] = health_C
	cmpnts["PlayerComponent"] = player_C
	cmpnts["InventoryComponent"] = inventory_C
	cmpnts["AttackComponent"] = attack_C
	cmpnts["AbilityComponent"] = ability_C
	cmpnts["PlayerLocationComponent"] = location_C
	cmpnts["StatsComponent"] = stats_C
	cmpnts["PlayerEquipmentComponent"] = equipment_C

	return cmpnts
}
