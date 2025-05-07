package service

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func makeRange(start, end int) []int {
	if start > end {
		return []int{}
	}

	a := make([]int, end-start)
	for i := range a {
		a[i] = start + i
	}

	return a
}

func doParallel[T any](values []T, fn func(value T) error) error {
	eg, _ := errgroup.WithContext(context.Background())

	for _, v := range values {
		eg.Go(func() error {
			return fn(v)
		})
	}

	return eg.Wait()
}
