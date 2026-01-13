package infra

import (
	"context"
	"net/http"
	"testing"
	"time"
	"wfs/backend/config"
	"wfs/backend/data"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClanClient_ClanAutoComplete(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		expected := data.ClanAutocomplete{
			SearchAutocompleteResult: []struct {
				HexColor string      `json:"hex_color"`
				Tag      string      `json:"tag"`
				ID       data.ClanID `json:"id"`
			}{
				{HexColor: "#000000", Tag: "TEST", ID: 0},
				{HexColor: "#000001", Tag: "TEST2", ID: 1},
			},
		}
		server := simpleMockServer(t, 200, expected)
		defer server.Close()

		injector := do.New()
		do.ProvideValue(injector, config.Config{
			ClanClient: config.ClanConfig{
				URL:        server.URL,
				RetryCount: 0,
				Timeout:    0,
			},
		})
		instance, err := NewClanClient(injector)
		require.NoError(t, err)
		result, err := instance.ClanAutoComplete("TEST")

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("異常系_エラーレスポンス", func(t *testing.T) {
		t.Parallel()

		body := `{
            "status": "error",
            "data": {
                "search": [
                    "Length must be between 2 and 70."
                ]
            }
        }`
		server := simpleMockServer(t, http.StatusConflict, body)
		defer server.Close()

		injector := do.New()
		do.ProvideValue(injector, config.Config{
			ClanClient: config.ClanConfig{
				URL:        server.URL,
				RetryCount: 0,
				Timeout:    0,
			},
		})
		instance, err := NewClanClient(injector)
		require.NoError(t, err)
		_, err = instance.ClanAutoComplete("")

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
			ClanClient: config.ClanConfig{
				URL:        server.URL,
				RetryCount: 0,
				Timeout:    timeout - 1,
			},
		})
		instance, err := NewClanClient(injector)
		require.NoError(t, err)
		_, err = instance.ClanAutoComplete("")

		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
}
