package infra

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"wfs/backend/domain"
)

func TestNumbers_ExpectedStats(t *testing.T) {
	t.Parallel()

	// モック
	mockResponse := []byte(`{
        "time": 1621699200,
        "data": {
            "1234": {"average_damage_dealt": 50000, "average_frags": 1.2, "win_rate": 52.3},
            "5678": {"average_damage_dealt": 60000, "average_frags": 1.5, "win_rate": 56.8},
            "9012": []
        }
    }`)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(mockResponse)
	}))
	defer server.Close()

	// テスト
	numbers := NewNumbers(RequestConfig{URL: server.URL})
	actual, err := numbers.ExpectedStats()
	assert.NoError(t, err)

	// アサーション
	expected := domain.NSExpectedStats{
		Time: 1621699200,
		Data: map[int]domain.NSExpectedStatsData{
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
}
