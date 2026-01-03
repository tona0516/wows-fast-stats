package infra

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func simpleMockServer[T any](t *testing.T, statusCode int, body T) *httptest.Server {
	var byteBody []byte
	if converted, ok := any(body).(string); ok {
		byteBody = []byte(converted)
	} else {
		var err error
		byteBody, err = json.Marshal(body)
		require.NoError(t, err)
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write(byteBody)
	}))
}

func timeoutMockServer[T any](t *testing.T, statusCode int, body T, timeout time.Duration) *httptest.Server {
	var byteBody []byte
	if converted, ok := any(body).(string); ok {
		byteBody = []byte(converted)
	} else {
		var err error
		byteBody, err = json.Marshal(body)
		require.NoError(t, err)
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(timeout)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write(byteBody)
	}))
}
