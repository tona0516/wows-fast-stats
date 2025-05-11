package clans

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

const provideName = "ClansAPIClient"

func TestClansWargaming_FetchAutoComplete(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		injector := do.New()
		do.ProvideNamed(injector, provideName, func(i *do.Injector) (*req.Client, error) {
			server := testutil.NewStubServer(t, http.StatusOK, map[string]interface{}{
				"search_autocomplete_result": []interface{}{
					map[string]interface{}{
						"tag":       "-K2-",
						"name":      "神風-s",
						"id":        2000036632,
						"hex_color": "#cc9966",
					},
				},
				"_meta_": map[string]interface{}{
					"collection":  "search_autocomplete_result",
					"total_clans": 1,
				},
			})

			client := req.C()
			client.SetBaseURL(server.URL)

			return client, nil
		})

		clansWargaming := NewAPI(injector)
		result, err := clansWargaming.FetchAutoComplete("-K2-")

		assert.NoError(t, err)
		assert.Equal(t, 1, len(result.SearchAutocompleteResult))
		assert.Equal(t, "-K2-", result.SearchAutocompleteResult[0].Tag)
		assert.Equal(t, "#cc9966", result.SearchAutocompleteResult[0].HexColor)
		assert.Equal(t, 2000036632, result.SearchAutocompleteResult[0].ID)
	})

	t.Run("異常系", func(t *testing.T) {
		t.Parallel()

		injector := do.New()
		do.ProvideNamed(injector, provideName, func(i *do.Injector) (*req.Client, error) {
			server := testutil.NewStubServer(t, http.StatusConflict, map[string]interface{}{
				"status": "error",
				"data": map[string]interface{}{
					"search": []interface{}{
						"Length must be between 2 and 70.",
					},
				},
			})

			client := req.C()
			client.SetBaseURL(server.URL)

			return client, nil
		})

		clansWargaming := NewAPI(injector)
		_, err := clansWargaming.FetchAutoComplete("a")

		assert.Error(t, err)
		assert.True(t, failure.Is(err, apperr.UWGAPIError))
	})
}
