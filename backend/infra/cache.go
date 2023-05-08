package infra

import (
	"changeme/backend/apperr"
	"encoding/gob"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const CacheDir string = "cache"

type Cache[T any] struct {
	Name string
}

func (c *Cache[T]) Serialize(target T) error {
	errDetail := apperr.Cache.Serialize

	_ = os.Mkdir(CacheDir, 0o755)

	f, err := os.Create(filepath.Join(CacheDir, c.Name+".bin"))
	if err != nil {
		return errors.WithStack(errDetail.WithRaw(err))
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	if err := enc.Encode(target); err != nil {
		return errors.WithStack(errDetail.WithRaw(err))
	}

	return nil
}

func (c *Cache[T]) Deserialize() (T, error) {
	errDetail := apperr.Cache.Serialize

	_ = os.Mkdir(CacheDir, 0o755)

	var object T
	f, err := os.Open(filepath.Join(CacheDir, c.Name+".bin"))
	if err != nil {
		return object, errors.WithStack(errDetail.WithRaw(err))
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	if err := dec.Decode(&object); err != nil {
		return object, errors.WithStack(errDetail.WithRaw(err))
	}

	return object, nil
}
