package data

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/morikuni/failure"
)

const (
	NumbersAvgDamage = "average_damage_dealt"
	NumbersAvgFrags  = "average_frags"
	NumbersWinrate   = "win_rate"
)

var ErrNSExpectedStatsNoDataKey = errors.New("NSExpectedStats: no data key")

type NSExpectedStats struct {
	Data NSExpectedStatsData `json:"data"`
}

type NSExpectedStatsData map[int]NSExpectedStatsValues

type NSExpectedStatsValues struct {
	AverageDamageDealt float64 `json:"average_damage_dealt"`
	AverageFrags       float64 `json:"average_frags"`
	WinRate            float64 `json:"win_rate"`
}

func (n *NSExpectedStats) UnmarshalJSON(b []byte) error {
	root := make(map[string]any)
	if err := json.Unmarshal(b, &root); err != nil {
		return failure.Wrap(err)
	}

	data, ok := root["data"].(map[string]any)
	if !ok {
		return failure.Wrap(ErrNSExpectedStatsNoDataKey)
	}

	es := make(NSExpectedStatsData)
	for key, value := range data {
		shipID, err := strconv.Atoi(key)
		if err != nil {
			continue
		}

		values, ok := value.(map[string]any)
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

		es[shipID] = NSExpectedStatsValues{
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
