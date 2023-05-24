package infra

import (
	"changeme/backend/apperr"
	"encoding/gob"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type Cache[T any] struct {
	Name string
	Dir  string
}

func NewCache[T any](name string, dir string) *Cache[T] {
	return &Cache[T]{
		Name: name,
		Dir:  dir,
	}
}

func (c *Cache[T]) Serialize(target T) error {
	e := apperr.Cache.Serialize

	_ = os.MkdirAll(c.Dir, 0o755)

	f, err := os.Create(filepath.Join(c.Dir, c.Name+".bin"))
	if err != nil {
		return errors.WithStack(e.WithRaw(err))
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	if err := enc.Encode(target); err != nil {
		return errors.WithStack(e.WithRaw(err))
	}

	return nil
}

func (c *Cache[T]) Deserialize() (T, error) {
	e := apperr.Cache.Deserialize

	_ = os.MkdirAll(c.Dir, 0o755)

	var object T
	f, err := os.Open(filepath.Join(c.Dir, c.Name+".bin"))
	if err != nil {
		return object, errors.WithStack(e.WithRaw(err))
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	if err := dec.Decode(&object); err != nil {
		return object, errors.WithStack(e.WithRaw(err))
	}

	return object, nil
}
