package infra

import (
	"encoding/gob"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const CACHE_DIRECTORY = "cache"

type Cache[T any] struct {
	Name string
}

func (c *Cache[T]) Serialize(object T) error {
	_ = os.Mkdir(CACHE_DIRECTORY, 0755)

    filename := c.Name + ".bin"
    path := filepath.Join(CACHE_DIRECTORY, filename)
	f, err := os.Create(path)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
    return errors.WithStack(enc.Encode(object))
}

func (c *Cache[T]) Deserialize() (T, error) {
	var object T

    filename := c.Name + ".bin"
    path := filepath.Join(CACHE_DIRECTORY, filename)
	f, err := os.Open(path)
	if err != nil {
		return object, errors.WithStack(err)
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	if err := dec.Decode(&object); err != nil {
		return object, errors.WithStack(err)
	}
	return object, nil
}
