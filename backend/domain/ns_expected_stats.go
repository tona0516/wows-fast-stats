package domain

import (
	"encoding/json"
	"strconv"
	"wfs/backend/apperr"

	"github.com/morikuni/failure"
)

const (
	NumbersAvgDamage = "average_damage_dealt"
	NumbersAvgFrags  = "average_frags"
	NumbersWinrate   = "win_rate"
)

type NSExpectedStats struct {
	Time int64                       `json:"time"`
	Data map[int]NSExpectedStatsData `json:"data"`
}

func (n *NSExpectedStats) UnmarshalJSON(b []byte) error {
	errCtx := failure.Context{"body": string(b)}

	root := make(map[string]interface{})
	if err := json.Unmarshal(b, &root); err != nil {
		return failure.New(apperr.ExpectedStatsParseError, errCtx, failure.Messagef("%s", err.Error()))
	}

	time, ok := root["time"].(float64)
	if !ok {
		return failure.New(apperr.ExpectedStatsParseError, errCtx, failure.Messagef("%s", "no time key"))
	}
	data, ok := root["data"].(map[string]interface{})
	if !ok {
		return failure.New(apperr.ExpectedStatsParseError, errCtx, failure.Messagef("%s", "no data key"))
	}

	esd := make(map[int]NSExpectedStatsData)
	for key, value := range data {
		shipID, err := strconv.Atoi(key)
		if err != nil {
			continue
		}

		values, ok := value.(map[string]interface{})
		if !ok {
			continue
		}

		damage, ok := values[NumbersAvgDamage].(float64)
		if !ok {
			continue
		}

		frags, ok := values[NumbersAvgFrags].(float64)
		if !ok {
			continue
		}

		wr, ok := values[NumbersWinrate].(float64)
		if !ok {
			continue
		}

		esd[shipID] = NSExpectedStatsData{
			AverageDamageDealt: damage,
			AverageFrags:       frags,
			WinRate:            wr,
		}
	}

	*n = NSExpectedStats{
		Time: int64(time),
		Data: esd,
	}

	return nil
}

type NSExpectedStatsData struct {
	AverageDamageDealt float64 `json:"average_damage_dealt"`
	AverageFrags       float64 `json:"average_frags"`
	WinRate            float64 `json:"win_rate"`
}
