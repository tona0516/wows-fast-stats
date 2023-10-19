package infra

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"wfs/backend/apperr"
	"wfs/backend/domain"
)

func TestNumbers_ExpectedStats(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		// モック
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
		assert.NoError(t, err)

		// アサーション
		expected := domain.NSExpectedStats{
			Time: 1621699200,
			Data: map[int]domain.ExpectedStats{
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
			},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("異常系", func(t *testing.T) {
		t.Parallel()

		server := simpleMockServer(500, `{
            "time": 0,
            "data": {}
        }`)
		defer server.Close()

		// テスト
		numbers := NewNumbers(RequestConfig{URL: server.URL})
		_, err := numbers.ExpectedStats()
		assert.EqualError(t, apperr.Unwrap(err), apperr.NumbersAPIFetchExpectedStatsError.ErrorCode())
	})
}
