package vo

import "reflect"

type WGAccountInfo struct {
	Status string                    `json:"status"`
	Data   map[int]WGAccountInfoData `json:"data"`
	Error  WGError                   `json:"error"`
}

func (w WGAccountInfo) GetStatus() string {
	return w.Status
}

func (w WGAccountInfo) GetError() WGError {
	return w.Error
}

type WGAccountInfoData struct {
	HiddenProfile bool `json:"hidden_profile"`
	Statistics    struct {
		Pvp struct {
			Wins            uint `json:"wins"`
			Battles         uint `json:"battles"`
			DamageDealt     uint `json:"damage_dealt"`
			Frags           uint `json:"frags"`
			SurviveWins     uint `json:"survived_wins"`
			SurvivedBattles uint `json:"survived_battles"`
			Xp              uint `json:"xp"`
		} `json:"pvp"`
	}
}

func (w WGAccountInfoData) Field() string {
	return fieldQuery(reflect.TypeOf(&w).Elem())
}
