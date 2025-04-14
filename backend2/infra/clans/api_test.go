package clans

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
	client := req.C().
		SetBaseURL(server.URL)

	return client
}

func TestClansWargaming_FetchAutoComplete(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		server := newMockServer(http.StatusOK, `{"search_autocomplete_result": [{"tag": "-K2-", "name": "\u795e\u98a8-s", "id": 2000036632, "hex_color": "#cc9966"}], "_meta_": {"collection": "search_autocomplete_result", "total_clans": 1}}`)
		defer server.Close()

		client := newMockClient(server)
		clansWargaming := NewAPI(client)

		result, err := clansWargaming.FetchAutoComplete("-K2-")

		assert.NoError(t, err)
		assert.Equal(t, 1, len(result.SearchAutocompleteResult))
		assert.Equal(t, "-K2-", result.SearchAutocompleteResult[0].Tag)
		assert.Equal(t, "#cc9966", result.SearchAutocompleteResult[0].HexColor)
		assert.Equal(t, 2000036632, result.SearchAutocompleteResult[0].ID)
	})

	t.Run("異常系", func(t *testing.T) {
		t.Parallel()

		server := newMockServer(http.StatusConflict, `{"status": "error", "data": {"search": ["Length must be between 2 and 70."]}}`)
		defer server.Close()

		client := newMockClient(server)
		clansWargaming := NewAPI(client)
		_, err := clansWargaming.FetchAutoComplete("")

		assert.Error(t, err)
		assert.True(t, failure.Is(err, apperr.UWGAPIError))
	})
}
