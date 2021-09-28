package misc

import (
	"encoding/json"
	"log"
)

const (
	GO_TO = "GT"

	//Loot actions
	PICK_UP_ITEM      = "PUI"
	PICK_UP_ALL_ITEMS = "PUAI"
	SELECT_LOOT_ITEM  = "SLI"
	DISMISS_LOOT      = "DL"

	//Expedition actions
	EXPEDITION_CONTINUE          = "EC"
	EXPEDITION_LEAVE             = "EL"
	EXPEDITION_ACTIVITY_SELECTED = "EAS"

	//Gathering Actions
	GATHERING_CORRECT   = "GC"
	GATHERING_INCORRECT = "GI"

	//Crafting table
	CRAFTING_TABLE_RECIPE_SELECTED = "CTRS"

	//Inventory window
	INVENTORY_ITEM_SELECTED = "IIS"
	INVENTORY_ITEM_USE      = "IIU"
	INVENTORY_DROP_ITEM     = "IDI"

	//Character window
	CHARACTER_SELECT_TAB         = "CST"
	CHARACTER_EQUIPMENT_SELECTED = "CES"

	//COMMON
	WINDOW_CLOSE = "WC"
)

//CallbackData - type used to wrap data for inline keyboard buttons
//Maximum length - 64b
type CallbackData struct {
	Action  string `json:"a"`
	ID      string `json:"id"`
	Payload string `json:"p"`
}

//JSON - converts `CallbackData` to JSON string
func (cb CallbackData) JSON() string {
	cbJSON, err := json.Marshal(cb)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(cbJSON)
}

//FromJSON - converts JSON string to `CallbackData`
func (cb *CallbackData) FromJSON(jsonStr string) {
	log.Println(jsonStr)
	err := json.Unmarshal([]byte(jsonStr), &cb)
	if err != nil {
		log.Println(err)
	}
}
