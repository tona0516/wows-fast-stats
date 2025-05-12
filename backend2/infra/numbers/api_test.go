package numbers

import (
	"net/http"
	"testing"
	"wfs/backend2/testutil"

	"github.com/imroc/req/v3"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
)

const provideName = "NumbersAPIClient"

func TestAPI_Fetch(t *testing.T) {
	t.Parallel()

	injector := do.New()
	do.ProvideNamed(injector, provideName, func(i *do.Injector) (*req.Client, error) {
		server := testutil.NewStubServer(t, http.StatusOK, map[string]interface{}{
			"time": 1621699200,
			"data": map[string]interface{}{
				"1234": map[string]interface{}{
					"average_damage_dealt": 50000,
					"average_frags":        1.2,
					"win_rate":             52.3,
				},
				"5678": map[string]interface{}{
					"average_damage_dealt": 60000,
					"average_frags":        1.5,
					"win_rate":             56.8,
				},
				"9012": []interface{}{},
			},
		})

		client := req.C()
		client.SetBaseURL(server.URL)

		return client, nil
	})

	api := NewAPI(injector)
	actual, err := api.Fetch()

	assert.NoError(t, err)

	expected := Expected{
		Data: map[int]ExpectedValues{
			1234: {
				AverageDamageDealt: 50000,
				AverageFrags:       1.2,
				WinRate:            52.3,
			},
			5678: {
				AverageDamageDealt: 60000,
				AverageFrags:       1.5,
				WinRate:            56.8,
			},
		},
	}
	assert.Equal(t, expected, actual)
}
