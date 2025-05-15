package common

import (
	"bytes"
	"context"
	"encoding/json"

	"golang.org/x/sync/errgroup"
)

func MakeRange(start, end int) []int {
	if start > end {
		return []int{}
	}

	a := make([]int, end-start)
	for i := range a {
		a[i] = start + i
	}

	return a
}

func DoParallel[T any](values []T, fn func(value T) error) error {
	eg, _ := errgroup.WithContext(context.Background())

	for _, v := range values {
		eg.Go(func() error {
			return fn(v)
		})
	}

	return eg.Wait()
}

func PrettyJSON(v any) string {
	jsonByte, err := json.Marshal(v)
	if err != nil {
		return ""
	}

	var buf bytes.Buffer
	if err := json.Indent(&buf, jsonByte, "", "    "); err != nil {
		return ""
	}

	return buf.String()
}
