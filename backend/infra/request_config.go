package infra

import (
	"time"
)

type RequestConfig struct {
	URL     string
	Retry   uint64
	Timeout time.Duration
}
