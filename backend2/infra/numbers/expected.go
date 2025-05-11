package numbers

import (
	"encoding/json"
	"strconv"
	"wfs/backend/apperr"

	"github.com/morikuni/failure"
)

type ExpectedValues struct {
	AverageDamageDealt float64 `json:"average_damage_dealt"`
	AverageFrags       float64 `json:"average_frags"`
	WinRate            float64 `json:"win_rate"`
}

type Expected struct {
	Data map[int]ExpectedValues `json:"data"`
}

func (e *Expected) UnmarshalJSON(b []byte) error {
	errCtx := failure.Context{"body": string(b)}

	root := make(map[string]interface{})
	if err := json.Unmarshal(b, &root); err != nil {
		return failure.Translate(err, apperr.ParseExpectedStatsError, errCtx)
	}

	data, ok := root["data"].(map[string]interface{})
	if !ok {
		return failure.New(apperr.ParseExpectedStatsError, errCtx)
	}

	es := make(map[int]ExpectedValues)
	for key, value := range data {
		shipID, err := strconv.Atoi(key)
		if err != nil {
			continue
		}

		values, ok := value.(map[string]interface{})
		if !ok {
			continue
		}

		damage, ok := values["average_damage_dealt"].(float64)
		if !ok {
			continue
		}

		frags, ok := values["average_frags"].(float64)
		if !ok {
			continue
		}

		wr, ok := values["win_rate"].(float64)
		if !ok {
			continue
		}

		es[shipID] = ExpectedValues{
			AverageDamageDealt: damage,
			AverageFrags:       frags,
			WinRate:            wr,
		}
	}

	*e = Expected{
		Data: es,
	}

	return nil
}
