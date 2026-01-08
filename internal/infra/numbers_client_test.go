package infra

import (
	"context"
	"net/http"
	"testing"
	"time"
	"wfs/internal/config"
	"wfs/internal/data"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNumbersClient_ExpectedStats(t *testing.T) {
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

		injector := do.New()
		do.ProvideValue(injector, config.Config{
			NumbersClient: config.NumbersClientConfig{
				URL:        server.URL,
				RetryCount: 0,
				Timeout:    0,
			},
		})
		instance, err := NewNumbersClient(injector)
		require.NoError(t, err)
		actual, err := instance.ExpectedStats()

		assert.NoError(t, err)
		assert.Equal(t, data.NSExpectedStats{
			Data: data.NSExpectedStatsData{
				1234: data.NSExpectedStatsValues{
					AverageDamageDealt: 50000,
					AverageFrags:       1.2,
					WinRate:            52.3,
				},
				5678: data.NSExpectedStatsValues{
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

		injector := do.New()
		do.ProvideValue(injector, config.Config{
			NumbersClient: config.NumbersClientConfig{
				URL:        server.URL,
				RetryCount: 0,
				Timeout:    0,
			},
		})
		instance, err := NewNumbersClient(injector)
		require.NoError(t, err)
		_, err = instance.ExpectedStats()

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

		injector := do.New()
		do.ProvideValue(injector, config.Config{
			NumbersClient: config.NumbersClientConfig{
				URL:        server.URL,
				RetryCount: 0,
				Timeout:    timeout - 1,
			},
		})
		instance, err := NewNumbersClient(injector)
		require.NoError(t, err)
		_, err = instance.ExpectedStats()

		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
}
