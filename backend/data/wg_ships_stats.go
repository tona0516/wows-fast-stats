package data

type AllPlayerShipsStats map[int]WGShipsStats

func (w AllPlayerShipsStats) Player(accountID int) []WGShipsStatsData {
	return w[accountID].Data[accountID]
}

type WGShipsStats struct {
	WGResponseCommon[map[int][]WGShipsStatsData]
}

type WGShipsStatsData struct {
	Pvp     WGShipStatsValues `json:"pvp"`
	PvpSolo WGShipStatsValues `json:"pvp_solo"`
	PvpDiv2 struct {
		Battles uint `json:"battles"`
	} `json:"pvp_div2"`
	PvpDiv3 struct {
		Battles uint `json:"battles"`
	} `json:"pvp_div3"`
	RankSolo WGShipStatsValues `json:"rank_solo"`
	ShipID   int               `json:"ship_id"`
}

type WGShipStatsValues struct {
	Wins            uint       `json:"wins"`
	Battles         uint       `json:"battles"`
	DamageDealt     uint       `json:"damage_dealt"`
	MaxDamageDealt  uint       `json:"max_damage_dealt"`
	Frags           uint       `json:"frags"`
	SurvivedWins    uint       `json:"survived_wins"`
	SurvivedBattles uint       `json:"survived_battles"`
	Xp              uint       `json:"xp"`
	MainBattery     WGArmament `json:"main_battery"`
	Torpedoes       WGArmament `json:"torpedoes"`
	PlanesKilled    uint       `json:"planes_killed"`
}

type WGArmament struct {
	Hits  uint `json:"hits"`
	Shots uint `json:"shots"`
}
