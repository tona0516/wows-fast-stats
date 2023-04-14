package vo

type WGAccountInfo struct {
	Status string `json:"status"`
	Meta   struct {
		Count  int         `json:"count"`
		Hidden interface{} `json:"hidden"`
	} `json:"meta"`
	Data map[int]struct {
        HiddenProfile bool `json:"hidden_profile"`
		Statistics struct {
			Pvp struct {
				Wins            uint `json:"wins"`
				Battles         uint `json:"battles"`
				DamageDealt     uint `json:"damage_dealt"`
				Frags           uint `json:"frags"`
				SurvivedBattles uint `json:"survived_battles"`
			} `json:"pvp"`
		} `json:"statistics"`
	} `json:"data"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Field   string `json:"field"`
		Value   string `json:"value"`
	} `json:"error"`
}
