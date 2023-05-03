package service

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

func makeRange(min, max int) []int {
    a := make([]int, max-min+1)
    for i := range a {
        a[i] = min + i
    }
    return a
}

func doParallel[T any](maxParallels uint, values []T, fn func(value T) error) error {
    wg := sync.WaitGroup{}
    eg, ctx := errgroup.WithContext(context.Background())
	limit := make(chan struct{}, maxParallels)

    for _, v := range values {
        v := v
        eg.Go(func() error {
            limit <- struct{}{}
            wg.Add(1)

            defer func() {
                wg.Done()
                <-limit
            }()

            select {
            case <-ctx.Done():
                return nil
            default:
                return fn(v)
            }
        })
    }

    return eg.Wait()
}
