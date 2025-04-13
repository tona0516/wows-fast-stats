package wargaming

import (
	"reflect"
)

type AccountInfoResponse struct {
	ResponseCommon[AccountInfo]
}

func (w AccountInfoResponse) Field() string {
	return fieldQuery(reflect.TypeOf(&AccountInfoData{}).Elem())
}

type AccountInfo map[int]AccountInfoData

type AccountInfoData struct {
	HiddenProfile bool `json:"hidden_profile"`
	Statistics    struct {
		Pvp      PlayerStatsValues `json:"pvp"`
		PvpSolo  PlayerStatsValues `json:"pvp_solo"`
		PvpDiv2  PlayerStatsValues `json:"pvp_div2"`
		PvpDiv3  PlayerStatsValues `json:"pvp_div3"`
		RankSolo PlayerStatsValues `json:"rank_solo"`
	} `json:"statistics"`
}

type PlayerStatsValues struct {
	Wins                 uint `json:"wins"`
	Battles              uint `json:"battles"`
	DamageDealt          uint `json:"damage_dealt"`
	MaxDamageDealt       uint `json:"max_damage_dealt"`
	MaxDamageDealtShipID int  `json:"max_damage_dealt_ship_id"`
	Frags                uint `json:"frags"`
	SurvivedWins         uint `json:"survived_wins"`
	SurvivedBattles      uint `json:"survived_battles"`
	Xp                   uint `json:"xp"`
}
