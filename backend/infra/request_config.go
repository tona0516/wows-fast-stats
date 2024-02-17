package infra

import (
	"net/http"
	"time"
)

type RequestConfig struct {
	URL       string
	Retry     uint64
	Timeout   time.Duration
	Transport *http.Transport
}
