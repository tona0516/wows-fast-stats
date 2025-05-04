package infra

import (
	"reflect"
)

type WGAccountInfoResponse struct {
	WGResponseCommon[WGAccountInfo]
}

func (w WGAccountInfoResponse) Field() string {
	return fieldQuery(reflect.TypeOf(&WGAccountInfoData{}).Elem())
}

type WGAccountInfo map[int]WGAccountInfoData

type WGAccountInfoData struct {
	HiddenProfile bool `json:"hidden_profile"`
	Statistics    struct {
		Pvp      WGPlayerStatsValues `json:"pvp"`
		PvpSolo  WGPlayerStatsValues `json:"pvp_solo"`
		PvpDiv2  WGPlayerStatsValues `json:"pvp_div2"`
		PvpDiv3  WGPlayerStatsValues `json:"pvp_div3"`
		RankSolo WGPlayerStatsValues `json:"rank_solo"`
	} `json:"statistics"`
}

type WGPlayerStatsValues struct {
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
