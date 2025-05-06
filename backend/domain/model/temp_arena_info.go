package model

import (
	"strings"
	"time"
)

type TempArenaInfo struct {
	Vehicles   []Vehicle `json:"vehicles"`
	DateTime   string    `json:"dateTime"`
	MapID      int       `json:"mapId"`
	MatchGroup string    `json:"matchGroup"`
	PlayerName string    `json:"playerName"`
}

func (t *TempArenaInfo) AccountNames() []string {
	accountNames := make([]string, 0)

	for _, v := range t.Vehicles {
		// Note: Bot name in corp or ramdom battle.
		if strings.HasPrefix(v.Name, ":") && strings.HasSuffix(v.Name, ":") {
			continue
		}

		// Note: Bot name in operation.
		if strings.HasPrefix(v.Name, "IDS_OP") {
			continue
		}

		accountNames = append(accountNames, v.Name)
	}

	return accountNames
}

func (t *TempArenaInfo) Unixtime() int64 {
	date, err := time.ParseInLocation("02.01.2006 15:04:05", t.DateTime, time.Local)
	if err != nil {
		return 0
	}

	return date.Unix()
}

type Vehicle struct {
	ShipID   int    `json:"shipId"`
	Relation int    `json:"relation"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
}

func (v *Vehicle) IsFriend() bool {
	return v.Relation == 0 || v.Relation == 1
}
