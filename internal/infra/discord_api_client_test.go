package infra

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscordClient_Comment(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		server := simpleMockServer(t, http.StatusOK, map[string]string{})
		defer server.Close()

		instance, err := NewDiscordClient(server.URL, 0, 0)
		require.NoError(t, err)
		err = instance.Comment("test message")

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

		instance, err := NewDiscordClient(server.URL, 0, 0)
		require.NoError(t, err)
		err = instance.Comment("test message")

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

		instance, err := NewDiscordClient(server.URL, 0, timeout-1)
		require.NoError(t, err)
		err = instance.Comment("test message")

		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
}
