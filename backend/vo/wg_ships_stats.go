package vo

type WGShipsStats struct {
	Status string `json:"status"`
	Meta   struct {
		Count  int         `json:"count"`
		Hidden interface{} `json:"hidden"`
	} `json:"meta"`
	Data map[int][]struct {
		Pvp struct {
			Wins            uint `json:"wins"`
			Battles         uint `json:"battles"`
			DamageDealt     uint `json:"damage_dealt"`
			Frags           uint `json:"frags"`
			SurviveWins     uint `json:"survived_wins"`
			SurvivedBattles uint `json:"survived_battles"`
			Xp              uint `json:"xp"`
			MainBattery     struct {
				Hits  uint `json:"hits"`
				Shots uint `json:"shots"`
			} `json:"main_battery"`
			Torpedoes struct {
				Hits  uint `json:"hits"`
				Shots uint `json:"shots"`
			} `json:"torpedoes"`
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

func (w WGShipsStats) GetStatus() string {
	return w.Status
}

func (w WGShipsStats) GetError() WGError {
	return WGError(w.Error)
}
