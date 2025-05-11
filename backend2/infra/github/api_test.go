package github

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"wfs/backend/apperr"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
)

func newMockServer(statusCode int, responseBody string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(responseBody))
	}))
}

func newMockClient(server *httptest.Server) *req.Client {
	client := req.C()
	client.SetBaseURL(server.URL)

	return client
}

func TestAPI_FetchLatestRelease(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		server := newMockServer(200, `{
			"tag_name": "1.0.0",
			"html_url": "https://github.com/tona0516/wows-fast-stats/releases/tag/1.0.0"
		}`)
		defer server.Close()

		mockClient := newMockClient(server)
		api := NewGithub(mockClient)
		actual, err := api.FetchLatestRelease()

		assert.NoError(t, err)
		assert.Equal(t, "1.0.0", actual.TagName)
		assert.Equal(t, "https://github.com/tona0516/wows-fast-stats/releases/tag/1.0.0", actual.HTMLURL)
	})

	t.Run("異常系", func(t *testing.T) {
		t.Parallel()

		server := newMockServer(http.StatusBadRequest, "{}")
		defer server.Close()

		mockClient := newMockClient(server)
		api := NewGithub(mockClient)
		_, err := api.FetchLatestRelease()

		assert.Error(t, err)
		assert.True(t, failure.Is(err, apperr.GithubAPICheckUpdateError))
	})
}
