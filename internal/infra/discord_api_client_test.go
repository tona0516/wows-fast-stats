package infra

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDiscordClient_Comment(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		server := simpleMockServer(t, http.StatusOK, map[string]string{})
		defer server.Close()

		instance := NewDiscordClient(apiConfig{
			url:        server.URL,
			retryCount: 0,
			timeout:    0,
		})
		err := instance.Comment("test message")

		assert.NoError(t, err)
	})

	t.Run("異常系_エラーレスポンス", func(t *testing.T) {
		t.Parallel()

		body := `{
            "code": 10015,
            "message": "Invalid Webhook Token"
        }`
		server := simpleMockServer(t, http.StatusUnauthorized, body)
		defer server.Close()

		instance := NewDiscordClient(apiConfig{
			url:        server.URL,
			retryCount: 0,
			timeout:    0,
		})
		err := instance.Comment("test message")

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

		instance := NewDiscordClient(apiConfig{
			url:        server.URL,
			retryCount: 0,
			timeout:    timeout - 1,
		})
		err := instance.Comment("test message")

		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
}
