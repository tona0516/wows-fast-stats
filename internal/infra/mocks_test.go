package infra

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func simpleMockServer[T any](statusCode int, body T) *httptest.Server {
	var byteBody []byte
	if converted, ok := any(body).(string); ok {
		byteBody = []byte(converted)
	} else {
		//nolint:errchkjson
		byteBody, _ = json.Marshal(body)
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write(byteBody)
	}))
}
