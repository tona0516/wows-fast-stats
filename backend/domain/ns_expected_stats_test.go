package domain

import (
	"encoding/json"
	"testing"
	"wfs/backend/apperr"

	"github.com/stretchr/testify/assert"
)

func TestNSExpectedStats_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"time": 1621845600,
			"data": {
				"1": {
					"average_damage_dealt": 50000,
					"average_frags": 1.5,
					"win_rate": 0.6
				},
				"2": {
					"average_damage_dealt": 60000,
					"average_frags": 2.0,
					"win_rate": 0.7
				},
                "3": [],
                "4": {
					"average_damage_dealt": "value_not_float",
					"average_frags": 1.0,
					"win_rate": 0.1
				},
                "key_not_int": {
					"average_damage_dealt": 10000,
					"average_frags": 1.0,
					"win_rate": 0.1
				}
			}
		}`

		var actual NSExpectedStats
		err := json.Unmarshal([]byte(jsonData), &actual)
		assert.NoError(t, err)
		assert.Equal(t, NSExpectedStats{
			Time: int64(1621845600),
			Data: AllExpectedStats{
				1: {
					AverageDamageDealt: 50000.0,
					AverageFrags:       1.5,
					WinRate:            0.6,
				},
				2: {
					AverageDamageDealt: 60000.0,
					AverageFrags:       2.0,
					WinRate:            0.7,
				},
			},
		}, actual)
	})

	t.Run("異常系", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"time": "",
			"data": []
        }`

		err := json.Unmarshal([]byte(jsonData), &NSExpectedStats{})
		assert.EqualError(t, apperr.Unwrap(err), apperr.ParseExpectedStatsError.ErrorCode())
	})
}
