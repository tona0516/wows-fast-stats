package vo

import (
	"strings"

	"github.com/samber/lo"
)

type TempArenaInfo struct {
	Vehicles   []Vehicle `json:"vehicles"`
	DateTime   string    `json:"dateTime"`
	MapID      int       `json:"mapId"`
	MatchGroup string    `json:"matchGroup"`
	PlayerName string    `json:"playerName"`
}

type Vehicle struct {
	ShipID   int    `json:"shipId"`
	Relation int    `json:"relation"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
}

func (t *TempArenaInfo) AccountNames() []string {
	accountNames := lo.FilterMap(t.Vehicles, func(vehicle Vehicle, _ int) (string, bool) {
		// Note: Bot name in corp or ramdom battle.
		if strings.HasPrefix(vehicle.Name, ":") && strings.HasSuffix(vehicle.Name, ":") {
			return "", false
		}

		// Note: Bot name in operation.
		if strings.HasPrefix(vehicle.Name, "IDS_OP") {
			return "", false
		}

		return vehicle.Name, true
	})

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
