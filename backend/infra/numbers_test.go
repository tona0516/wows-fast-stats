package infra

import (
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/data"

	"github.com/stretchr/testify/assert"
)

func TestNumbers_ExpectedStats(t *testing.T) {
	t.Parallel()
	t.Run("正常系", func(t *testing.T) {
		t.Parallel()
		// 準備
		mockResponse := `{
        "time": 1621699200,
        "data": {
            "1234": {"average_damage_dealt": 50000, "average_frags": 1.2, "win_rate": 52.3},
            "5678": {"average_damage_dealt": 60000, "average_frags": 1.5, "win_rate": 56.8},
            "9012": []
        }
    }`
		server := simpleMockServer(200, mockResponse)
		defer server.Close()

		// テスト
		numbers := NewNumbers(RequestConfig{URL: server.URL})
		actual, err := numbers.ExpectedStats()

		// アサーション
		assert.NoError(t, err)
		expected := data.ExpectedStats{
			1234: {
				AverageDamageDealt: 50000,
				AverageFrags:       1.2,
				WinRate:            52.3,
			},
			5678: {
				AverageDamageDealt: 60000,
				AverageFrags:       1.5,
				WinRate:            56.8,
			},
		}
		assert.Equal(t, expected, actual)
	})
	t.Run("異常系", func(t *testing.T) {
		t.Parallel()

		// 準備
		server := simpleMockServer(500, `{
            "time": 0,
            "data": {}
        }`)
		defer server.Close()

		// テスト
		numbers := NewNumbers(RequestConfig{URL: server.URL})
		_, err := numbers.ExpectedStats()

		// アサーション
		assert.EqualError(t, apperr.Unwrap(err), apperr.NumbersAPIError.ErrorCode())
	})
}
