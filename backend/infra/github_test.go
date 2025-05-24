package infra

import (
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/data"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGithub_LatestRelease(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		expected := data.GHLatestRelease{
			HTMLURL: "https://github.com/tona0516/wows-fast-stats/releases/tag/1.0.0",
			TagName: "1.0.0",
		}
		server := simpleMockServer(200, expected)
		defer server.Close()

		github := NewGithub(server.URL, 0, 0)

		actual, err := github.LatestRelease()
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("異常系", func(t *testing.T) {
		t.Parallel()

		server := simpleMockServer(400, "{}")
		defer server.Close()

		github := NewGithub(server.URL, 0, 0)

		_, err := github.LatestRelease()
		require.EqualError(t, apperr.Unwrap(err), apperr.GithubAPICheckUpdateError.ErrorCode())
	})
}
