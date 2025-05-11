package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewStubServer(t *testing.T, statusCode int, body map[string]any) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bytes, err := json.Marshal(body)
		if err != nil {
			assert.Fail(t, err.Error())

			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write(bytes)
	}))
}
