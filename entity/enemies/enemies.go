package enemies

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"sync"

	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
)

var once sync.Once

type EnemiesCollection struct {
	rawEnemies map[int]interface{}
	manager    *entity.Manager
}

var instance *EnemiesCollection

func GetInstance() *EnemiesCollection {
	once.Do(func() {
		instance = &EnemiesCollection{}
		file, _ := ioutil.ReadFile("./resources/enemies.json")
		var enms interface{}
		err := json.Unmarshal([]byte(file), &enms)
		if err != nil {
			log.Println(err)
		}
		instance.rawEnemies = make(map[int]interface{})

		for _, enm := range enms.([]interface{}) {
			enemy_id := int(enm.(map[string]interface{})["id"].(float64))
			instance.rawEnemies[enemy_id] = enm
		}
	})

	return instance
}

func (ic *EnemiesCollection) SetManager(manager *entity.Manager) {
	ic.manager = manager
}

func (ic *EnemiesCollection) GetEnemyById(id int) entity.Entity {

	rawEnemy := ic.rawEnemies[id]
	enemy := parseEntity(rawEnemy, ic.manager)

	return *enemy
}

/* Manually map fields to the enemy */
func parseEntity(raw interface{}, manager *entity.Manager) *entity.Entity {
	log.Println("CREATING ENEMY")
	enemy := manager.NewEntity()
	for key, element := range raw.(map[string]interface{})["components"].(map[string]interface{}) {
		// To benefit from built-in parses, primitive.M component is converted
		// to bson and later converted into struct of the proper component type
		jsonElement, _ := json.Marshal(element)

		switch key {
		case "AttackComponent":
			var comp components.AttackComponent
			json.Unmarshal(jsonElement, &comp)
			enemy.AddComponent(&comp)
		case "HealthComponent":
			var comp components.HealthComponent
			json.Unmarshal(jsonElement, &comp)
			enemy.AddComponent(&comp)
		case "NameComponent":
			var comp components.NameComponent
			json.Unmarshal(jsonElement, &comp)
			enemy.AddComponent(&comp)
		case "NPCComponent":
			var comp components.NPCComponent
			json.Unmarshal(jsonElement, &comp)
			// ldc := components.LootDropConfig{Chance: 0.1, PossibleLoot: items.ItemBundle{ID: 1, Qty: 3}}
			// comp = components.NPCComponent{
			// 	PossibleLoot: []components.LootDropConfig{ldc},
			// }
			for i, cfg := range element.(map[string]interface{})["PossibleLoot"].([]interface{}) {
				min := int(cfg.(map[string]interface{})["possibleLoot"].(map[string]interface{})["qty_min"].(float64))
				max := int(cfg.(map[string]interface{})["possibleLoot"].(map[string]interface{})["qty_max"].(float64))
				comp.PossibleLoot[i].PossibleLoot.Qty = rand.Intn(max-min) + min
			}
			enemy.AddComponent(&comp)
		}
		enemy.AddComponent(&components.EffectsComponent{})
	}
	return enemy

}
