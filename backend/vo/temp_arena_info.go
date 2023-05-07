package vo

import (
	"strings"
)

type TempArenaInfo struct {
	Vehicles []struct {
		ShipID   int    `json:"shipId"`
		Relation int    `json:"relation"`
		ID       int    `json:"id"`
		Name     string `json:"name"`
	} `json:"vehicles"`
	DateTime   string `json:"dateTime"`
	MapID      int    `json:"mapId"`
	MatchGroup string `json:"matchGroup"`
	PlayerName string `json:"playerName"`
}

func (t *TempArenaInfo) AccountNames() []string {
	accountNames := make([]string, 0)
	for i := range t.Vehicles {
		vehicle := t.Vehicles[i]
		if strings.HasPrefix(vehicle.Name, ":") && strings.HasSuffix(vehicle.Name, ":") {
			continue
		}

		accountNames = append(accountNames, vehicle.Name)
	}

	return accountNames
}

func (t *TempArenaInfo) FormattedDateTime() string {
	datetimeArray := strings.Split(t.DateTime, " ")
	if len(datetimeArray) < 2 {
		return ""
	}
	dateArray := strings.Split(datetimeArray[0], ".")
	if len(dateArray) < 3 {
		return ""
	}

	return dateArray[2] + "-" + dateArray[1] + "-" + dateArray[0] + " " + datetimeArray[1]
}

func (t *TempArenaInfo) BattleArena(battleArenas WGBattleArenas) string {
	return battleArenas.Data[t.MapID].Name
}

func (t *TempArenaInfo) BattleType(battleTypes WGBattleTypes) string {
	rawBattleType := battleTypes.Data[strings.ToUpper(t.MatchGroup)].Name

	return strings.ReplaceAll(rawBattleType, " ", "")
}
