package usecase

import (
	"log"
	"time"
)

func measure(label string, fn func()) {
	start := time.Now()
	fn()
	log.Printf("%s: %dms\n", label, time.Since(start).Milliseconds())
}
