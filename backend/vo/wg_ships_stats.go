package vo

type WGShipsStats struct {
	Status string `json:"status"`
	Meta   struct {
		Count  int         `json:"count"`
		Hidden interface{} `json:"hidden"`
	} `json:"meta"`
	Data map[int][]struct {
		Pvp struct {
			Wins            int `json:"wins"`
			Battles         int `json:"battles"`
			DamageDealt     int `json:"damage_dealt"`
			Xp              int `json:"xp"`
			Frags           int `json:"frags"`
			SurvivedBattles int `json:"survived_battles"`
		} `json:"pvp"`
		ShipID int `json:"ship_id"`
	} `json:"data"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Field   string `json:"field"`
		Value   string `json:"value"`
	} `json:"error"`
}
