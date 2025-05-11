package github

import (
	"net/http"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend2/testutil"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
)

const provideName = "GithubAPIClient"

func TestAPI_FetchLatestRelease(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		injector := do.New()
		do.ProvideNamed(injector, provideName, func(i *do.Injector) (*req.Client, error) {
			server := testutil.NewStubServer(t, http.StatusOK, map[string]interface{}{
				"tag_name": "1.0.0",
				"html_url": "https://github.com/tona0516/wows-fast-stats/releases/tag/1.0.0",
			})

			client := req.C()
			client.SetBaseURL(server.URL)

			return client, nil
		})

		api := NewGithub(injector)
		actual, err := api.FetchLatestRelease()

		assert.NoError(t, err)
		assert.Equal(t, "1.0.0", actual.TagName)
		assert.Equal(t, "https://github.com/tona0516/wows-fast-stats/releases/tag/1.0.0", actual.HTMLURL)
	})

	t.Run("異常系", func(t *testing.T) {
		t.Parallel()

		injector := do.New()
		do.ProvideNamed(injector, provideName, func(i *do.Injector) (*req.Client, error) {
			server := testutil.NewStubServer(t, http.StatusBadRequest, map[string]interface{}{})

			client := req.C()
			client.SetBaseURL(server.URL)

			return client, nil
		})

		api := NewGithub(injector)
		_, err := api.FetchLatestRelease()

		assert.Error(t, err)
		assert.True(t, failure.Is(err, apperr.GithubAPICheckUpdateError))
	})
}
