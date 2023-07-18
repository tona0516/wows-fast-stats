package infra

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"wfs/backend/domain"

	"github.com/stretchr/testify/assert"
)

func TestGithub_LatestRelease(t *testing.T) {
	t.Parallel()

	expected := domain.GHLatestRelease{
		HTMLURL: "https://github.com/tona0516/wows-fast-stats/releases/tag/1.0.0",
		TagName: "1.0.0",
	}
	body, err := json.Marshal(expected)
	assert.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}))
	defer server.Close()

	github := NewGithub(RequestConfig{
		URL: server.URL,
	})

	actual, err := github.LatestRelease()

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
