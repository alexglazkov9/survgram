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
	//Activity actions
	ACTIVITY_CONTINUE = "AN"
	ACTIVITY_LEAVE    = "AL"
	// ACTION_JOIN              = "join"
	// ACTION_ATTACK            = "attack"
	// ACTION_PRIMARY_ABILITY   = "pr_a"
	// ACTION_SECONDARY_ABILITY = "se_a"
)

//CallbackData - type used to wrap data for inline keyboard buttons
//Maximum length - 64b
type CallbackData struct {
	Action  string `json:"a"`
	ID      int    `json:"id"`
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
	err := json.Unmarshal([]byte(jsonStr), &cb)
	if err != nil {
		log.Println(err)
	}
}
