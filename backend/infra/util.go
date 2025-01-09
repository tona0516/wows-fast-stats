package infra

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"wfs/backend/apperr"

	"github.com/morikuni/failure"
	"golang.org/x/sync/errgroup"
)

func readJSON[T any](path string, defaulValue T) (T, error) {
	errCtx := failure.Context{"path": path}

	f, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return defaulValue, failure.New(apperr.FileNotExist, errCtx)
		}
		return defaulValue, failure.Wrap(err, errCtx)
	}
	errCtx["target"] = string(f)

	err = json.Unmarshal(f, &defaulValue)
	return defaulValue, failure.Wrap(err, errCtx)
}

func writeJSON[T any](path string, target T) error {
	//nolint:errchkjson
	marshaled, _ := json.Marshal(target)
	errCtx := failure.Context{"path": path, "target": string(marshaled)}

	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	f, err := os.Create(path)
	if err != nil {
		return failure.Wrap(err, errCtx)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(target)
	return failure.Wrap(err, errCtx)
}

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

func doParallel[T any](values []T, fn func(value T) error) error {
	eg, _ := errgroup.WithContext(context.Background())

	for _, v := range values {
		eg.Go(func() error {
			return fn(v)
		})
	}

	return eg.Wait()
}
