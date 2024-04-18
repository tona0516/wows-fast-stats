package data

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

type ExpectedValues struct {
	AverageDamageDealt float64 `json:"average_damage_dealt"`
	AverageFrags       float64 `json:"average_frags"`
	WinRate            float64 `json:"win_rate"`
}

type ExpectedStats map[int]ExpectedValues

type NSExpectedStats struct {
	Data ExpectedStats `json:"data"`
}

func (n *NSExpectedStats) UnmarshalJSON(b []byte) error {
	errCtx := failure.Context{"body": string(b)}

	root := make(map[string]interface{})
	if err := json.Unmarshal(b, &root); err != nil {
		return failure.New(apperr.ParseExpectedStatsError, errCtx, failure.Messagef("%s", err.Error()))
	}

	data, ok := root["data"].(map[string]interface{})
	if !ok {
		return failure.New(apperr.ParseExpectedStatsError, errCtx, failure.Messagef("%s", "no data key"))
	}

	es := make(ExpectedStats)
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

		es[shipID] = ExpectedValues{
			AverageDamageDealt: damage,
			AverageFrags:       frags,
			WinRate:            wr,
		}
	}

	*n = NSExpectedStats{
		Data: es,
	}

	return nil
}
