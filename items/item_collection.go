package items

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

var once sync.Once

type ItemCollection struct {
	items map[int]IItem
}

var instance *ItemCollection

func GetItemCollection() *ItemCollection {
	start := time.Now()
	once.Do(func() {
		instance = &ItemCollection{}
		file, _ := ioutil.ReadFile("./resources/items.json")
		var itms interface{}
		_ = json.Unmarshal([]byte(file), &itms)
		instance.items = make(map[int]IItem)

		var item_struct IItem
		for _, itm := range itms.([]interface{}) {
			switch itm.(map[string]interface{})["type"].(string) {
			case string(WEAPON):
				item_struct = &Weapon{
					BaseItem: BaseItem{
						ID:     int(itm.(map[string]interface{})["id"].(float64)),
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
			case string(RESOURCE):
				item_struct = &Resource{
					BaseItem: BaseItem{
						ID:     int(itm.(map[string]interface{})["id"].(float64)),
						Name:   itm.(map[string]interface{})["name"].(string),
						Type:   ItemType(itm.(map[string]interface{})["type"].(string)),
						Tier:   int(itm.(map[string]interface{})["tier"].(float64)),
						Rarity: Rarity(itm.(map[string]interface{})["rarity"].(string)),
					},
				}
			case string(RECIPE):
				jsonElement, _ := json.Marshal(itm.(map[string]interface{})["ingridients"])
				var ingrs []ItemBundle
				json.Unmarshal(jsonElement, &ingrs)

				item_struct = &Recipe{
					BaseItem: BaseItem{
						ID:     int(itm.(map[string]interface{})["id"].(float64)),
						Name:   itm.(map[string]interface{})["name"].(string),
						Type:   ItemType(itm.(map[string]interface{})["type"].(string)),
						Tier:   int(itm.(map[string]interface{})["tier"].(float64)),
						Rarity: Rarity(itm.(map[string]interface{})["rarity"].(string)),
					},
					// Ingridients: ingrs,
					// Output:      int(itm.(map[string]interface{})["output"].(float64)),
				}
			case string(EQUIPMENT):
				item_struct = &Equipment{
					BaseItem: BaseItem{
						ID:     int(itm.(map[string]interface{})["id"].(float64)),
						Name:   itm.(map[string]interface{})["name"].(string),
						Type:   ItemType(itm.(map[string]interface{})["type"].(string)),
						Tier:   int(itm.(map[string]interface{})["tier"].(float64)),
						Rarity: Rarity(itm.(map[string]interface{})["rarity"].(string)),
					},
				}
			case string(TOOL):
				item_struct = &Tool{
					BaseItem: BaseItem{
						ID:     int(itm.(map[string]interface{})["id"].(float64)),
						Name:   itm.(map[string]interface{})["name"].(string),
						Type:   ItemType(itm.(map[string]interface{})["type"].(string)),
						Tier:   int(itm.(map[string]interface{})["tier"].(float64)),
						Rarity: Rarity(itm.(map[string]interface{})["rarity"].(string)),
					},
				}
			case string(BACKPACK):
				item_struct = &Backpack{
					BaseItem: BaseItem{
						ID:     int(itm.(map[string]interface{})["id"].(float64)),
						Name:   itm.(map[string]interface{})["name"].(string),
						Type:   ItemType(itm.(map[string]interface{})["type"].(string)),
						Tier:   int(itm.(map[string]interface{})["tier"].(float64)),
						Rarity: Rarity(itm.(map[string]interface{})["rarity"].(string)),
					},
				}
			case string(CONSUMABLE):
				item_struct = &Consumable{
					BaseItem: BaseItem{
						ID:     int(itm.(map[string]interface{})["id"].(float64)),
						Name:   itm.(map[string]interface{})["name"].(string),
						Type:   ItemType(itm.(map[string]interface{})["type"].(string)),
						Tier:   int(itm.(map[string]interface{})["tier"].(float64)),
						Rarity: Rarity(itm.(map[string]interface{})["rarity"].(string)),
					},
				}
			case string(IMPROVER):
				item_struct = &Improver{
					BaseItem: BaseItem{
						ID:     int(itm.(map[string]interface{})["id"].(float64)),
						Name:   itm.(map[string]interface{})["name"].(string),
						Type:   ItemType(itm.(map[string]interface{})["type"].(string)),
						Tier:   int(itm.(map[string]interface{})["tier"].(float64)),
						Rarity: Rarity(itm.(map[string]interface{})["rarity"].(string)),
					},
				}
			}
			instance.items[item_struct.GetID()] = item_struct
		}
	})
	elapsed := time.Since(start)
	log.Printf("Item initialization took %s for %d items", elapsed, len(instance.items))
	return instance
}

func (ic *ItemCollection) GetItemById(id int) IItem {
	item := ic.items[id]
	return item
}
