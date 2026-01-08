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

func TestGithubClient_LatestRelease(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		expected := data.GHLatestRelease{
			TagName: "1.0.0",
			HTMLURL: "https://github.com/tona0516/wows-fast-stats/releases/tag/1.0.0",
		}
		server := simpleMockServer(t, http.StatusOK, expected)
		defer server.Close()

		injector := do.New()
		do.ProvideValue(injector, config.Config{
			GithubClient: config.GithubClientConfig{
				URL:        server.URL,
				RetryCount: 0,
				Timeout:    0,
			},
		})
		instance, err := NewGithubClient(injector)
		require.NoError(t, err)
		result, err := instance.LatestRelease()

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("異常系_エラーレスポンス", func(t *testing.T) {
		t.Parallel()

		body := `{
            "message": "Not Found",
            "documentation_url": "https://docs.github.com/rest"
        }`
		server := simpleMockServer(t, http.StatusNotFound, body)
		defer server.Close()

		injector := do.New()
		do.ProvideValue(injector, config.Config{
			GithubClient: config.GithubClientConfig{
				URL:        server.URL,
				RetryCount: 0,
				Timeout:    0,
			},
		})
		instance, err := NewGithubClient(injector)
		require.NoError(t, err)
		_, err = instance.LatestRelease()

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
			GithubClient: config.GithubClientConfig{
				URL:        server.URL,
				RetryCount: 0,
				Timeout:    timeout - 1,
			},
		})
		instance, err := NewGithubClient(injector)
		require.NoError(t, err)
		_, err = instance.LatestRelease()

		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
}
