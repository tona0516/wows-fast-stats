package infra

import (
	"context"
	"net/http"
	"testing"
	"time"
	"wfs/internal/data"

	"github.com/stretchr/testify/assert"
)

func TestNumbersApiClient_ExpectedStats(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		body := `{
        	"time": 1621699200,
        	"data": {
            	"1234": {
					"average_damage_dealt": 50000,
					"average_frags": 1.2,
					"win_rate": 52.3
				},
            	"5678": {
					"average_damage_dealt": 60000,
					"average_frags": 1.5,
					"win_rate": 56.8
				},
            	"9012": []
        	}
    	}`
		server := simpleMockServer(t, http.StatusOK, body)
		defer server.Close()

		instance := NewNumbersApiClient(apiConfig{
			url:        server.URL,
			retryCount: 0,
			timeout:    0,
		})
		actual, err := instance.ExpectedStats()

		assert.NoError(t, err)
		assert.Equal(t, data.NSExpectedStats{
			Data: data.ExpectedStats{
				1234: data.ExpectedValues{
					AverageDamageDealt: 50000,
					AverageFrags:       1.2,
					WinRate:            52.3,
				},
				5678: data.ExpectedValues{
					AverageDamageDealt: 60000,
					AverageFrags:       1.5,
					WinRate:            56.8,
				},
			},
		}, actual)
	})

	t.Run("異常系_エラーレスポンス", func(t *testing.T) {
		t.Parallel()

		body := `{
            "time": 0,
            "data": {}
        }`
		server := simpleMockServer(t, http.StatusInternalServerError, body)
		defer server.Close()

		instance := NewNumbersApiClient(apiConfig{
			url:        server.URL,
			retryCount: 0,
			timeout:    0,
		})
		_, err := instance.ExpectedStats()

		assert.Error(t, err, ErrErrorResponse)
	})

	t.Run("異常系_タイムアウト", func(t *testing.T) {
		t.Parallel()

		timeout := time.Millisecond * 10
		server := timeoutMockServer(
			t,
			http.StatusOK,
			map[string]any{},
			timeout,
		)
		defer server.Close()

		instance := NewNumbersApiClient(apiConfig{
			url:        server.URL,
			retryCount: 0,
			timeout:    timeout - 1,
		})
		_, err := instance.ExpectedStats()

		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
}
