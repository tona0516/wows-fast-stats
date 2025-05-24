package infra

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/data"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnofficialWargaming_AccountListForSearch(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		expected := data.UWGClansAutocomplete{
			SearchAutocompleteResult: []struct {
				HexColor string `json:"hex_color"`
				Tag      string `json:"tag"`
				ID       int    `json:"id"`
			}{
				{HexColor: "#000000", Tag: "TEST", ID: 0},
				{HexColor: "#000001", Tag: "TEST2", ID: 1},
			},
		}

		server := simpleMockServer(200, expected)
		defer server.Close()

		uwargaming := NewUnofficialWargaming(server.URL, 0, 0)

		result, err := uwargaming.ClansAutoComplete("TEST")

		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("異常系", func(t *testing.T) {
		t.Parallel()

		body := `{
            "status": "error",
            "data": {
                "search": [
                    "Length must be between 2 and 70."
                ]
            }
        }`

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(body))
		}))
		defer server.Close()

		uwargaming := NewUnofficialWargaming(server.URL, 0, 0)

		_, err := uwargaming.ClansAutoComplete("")

		require.Error(t, err)
		code, ok := failure.CodeOf(err)
		assert.True(t, ok)
		assert.Equal(t, apperr.UWGAPIError, code)
	})
}
