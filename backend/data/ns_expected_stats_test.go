package data

import (
	"encoding/json"
	"testing"
	"wfs/backend/apperr"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNSExpectedStats_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		input := `{
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
                "5": {
                    "average_damage_dealt": 60000,
                    "average_frags": "value_not_float",
                    "win_rate": 0.1
                },
                "6": {
                    "average_damage_dealt": 60000,
                    "average_frags": 1.0,
                    "win_rate": "value_not_float"
                },
                "key_not_int": {
                    "average_damage_dealt": 10000,
                    "average_frags": 1.0,
                    "win_rate": 0.1
                }
            }
        }`

		var actual NSExpectedStats
		err := json.Unmarshal([]byte(input), &actual)
		require.NoError(t, err)
		assert.Equal(t, NSExpectedStats{
			Data: ExpectedStats{
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

		inputs := []string{
			"{}",
			`{
                "time": 1,
                "data": []
            }`,
		}

		for _, input := range inputs {
			err := json.Unmarshal([]byte(input), &NSExpectedStats{})
			require.Error(t, err)
			code, ok := failure.CodeOf(err)
			require.True(t, ok)
			// assert.Equal(t, apperr.ParseExpectedStatsError, code, fmt.Sprintf("actual=%s", code))
			assert.Equal(t, apperr.ParseExpectedStatsError, code)
		}
	})
}
