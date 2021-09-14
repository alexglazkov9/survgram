package enemies

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"

	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/components"
)

var once sync.Once

type EnemiesCollection struct {
	enemies map[int]*entity.Entity
}

var instance *EnemiesCollection

func GetInstance() *EnemiesCollection {
	once.Do(func() {
		instance = &EnemiesCollection{}
		file, _ := ioutil.ReadFile("./resources/enemies.json")
		var enms interface{}
		err := json.Unmarshal([]byte(file), &enms)
		log.Print("here")
		//log.Println(enms)
		if err != nil {
			log.Println(err)
		}
		instance.enemies = make(map[int]*entity.Entity)

		for _, enm := range enms.([]interface{}) {
			enemy_id := int(enm.(map[string]interface{})["id"].(float64))
			enemy := parseEntity(enm)
			instance.enemies[enemy_id] = enemy
		}
	})

	return instance
}

func (ic *EnemiesCollection) GetEnemyById(id int) entity.Entity {
	enemy := ic.enemies[id].Clone()

	return *enemy
}

/* Manually map fields to the enemy */
func parseEntity(raw interface{}) *entity.Entity {
	enemy := entity.New()
	for key, element := range raw.(map[string]interface{})["components"].(map[string]interface{}) {
		// To benefit from built-in parses, primitive.M component is converted
		// to bson and later converted into struct of the proper component type
		jsonElement, _ := json.Marshal(element)

		switch key {
		case "AttackComponent":
			var comp components.AttackComponent
			json.Unmarshal(jsonElement, &comp)
			log.Println(comp)
			enemy.AddComponent(&comp)
		case "HealthComponent":
			var comp components.HealthComponent
			json.Unmarshal(jsonElement, &comp)
			log.Println(comp)
			enemy.AddComponent(&comp)
		case "NameComponent":
			var comp components.NameComponent
			json.Unmarshal(jsonElement, &comp)
			log.Println(comp)
			enemy.AddComponent(&comp)
		}
		enemy.AddComponent(&components.EffectsComponent{})
	}
	return enemy

}
