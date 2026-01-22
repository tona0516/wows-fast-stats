package usecase

import (
	"time"
)

func measure(fn func()) int64 {
	start := time.Now()
	fn()
	return time.Since(start).Milliseconds()
}
