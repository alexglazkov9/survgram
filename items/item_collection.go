package items

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

var once sync.Once

type ItemCollection struct {
	items map[int]IItem
}

var instance *ItemCollection

func GetInstance() *ItemCollection {
	once.Do(func() {
		instance = &ItemCollection{}
		file, _ := ioutil.ReadFile("./resources/items.json")
		var itms interface{}
		_ = json.Unmarshal([]byte(file), &itms)
		instance.items = make(map[int]IItem)

		var item_struct IItem
		for _, itm := range itms.([]interface{}) {
			switch itm.(map[string]interface{})["type"].(string) {
			case string(MELEE):
				log.Println(itm.(map[string]interface{})["name"].(string))
				item_struct = &Weapon{
					BaseItem: BaseItem{
						ID:     int(itm.(map[string]interface{})["item_id"].(float64)),
						Name:   itm.(map[string]interface{})["name"].(string),
						Type:   ItemType(itm.(map[string]interface{})["type"].(string)),
						Tier:   int(itm.(map[string]interface{})["tier"].(float64)),
						Rarity: Rarity(itm.(map[string]interface{})["rarity"].(string)),
					},
					BaseDamage:      int(itm.(map[string]interface{})["damage"].(float64)),
					BaseAttackSpeed: itm.(map[string]interface{})["base_attack_speed"].(float64),
					Intellect:       int(itm.(map[string]interface{})["int"].(float64)),
					Strength:        int(itm.(map[string]interface{})["str"].(float64)),
					Agility:         int(itm.(map[string]interface{})["agi"].(float64)),
					MinLvlReq:       int(itm.(map[string]interface{})["min_lvl_req"].(float64)),
				}
			}
			instance.items[item_struct.GetID()] = item_struct
		}
		//log.Println(instance.items)
	})

	return instance
}

func (ic *ItemCollection) GetItemById(id int) IItem {
	item := ic.items[id]
	return item
}
