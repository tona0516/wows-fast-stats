package numbers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imroc/req/v3"
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

func TestAPI_Fetch(t *testing.T) {
	t.Parallel()

	server := newMockServer(200, `{
		"time": 1621699200,
		"data": {
			"1234": {"average_damage_dealt": 50000, "average_frags": 1.2, "win_rate": 52.3},
			"5678": {"average_damage_dealt": 60000, "average_frags": 1.5, "win_rate": 56.8},
			"9012": []
		}
	}`)
	defer server.Close()

	mockClient := newMockClient(server)
	api := NewAPI(mockClient)
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
