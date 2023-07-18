package service

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func makeRange(min, max int) []int {
	if min > max {
		return []int{}
	}

	a := make([]int, max-min)
	for i := range a {
		a[i] = min + i
	}

	return a
}

func doParallel[T any](limit uint, values []T, fn func(value T) error) error {
	eg, _ := errgroup.WithContext(context.Background())
	eg.SetLimit(int(limit))

	for _, v := range values {
		v := v
		eg.Go(func() error {
			return fn(v)
		})
	}

	return eg.Wait()
}
