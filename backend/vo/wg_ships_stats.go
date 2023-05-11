package vo

import "reflect"

type WGShipsStats struct {
	Status string                     `json:"status"`
	Data   map[int][]WGShipsStatsData `json:"data"`
	Error  WGError                    `json:"error"`
}

func (w WGShipsStats) GetStatus() string {
	return w.Status
}

func (w WGShipsStats) GetError() WGError {
	return w.Error
}

type WGShipsStatsData struct {
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
}

func (w WGShipsStatsData) Field() string {
	return fieldQuery(reflect.TypeOf(&w).Elem())
}
