package infra

import (
	"changeme/backend/apperr"
	"encoding/gob"
	"os"
	"path/filepath"

	"github.com/morikuni/failure"
)

const CacheDir string = "cache"

type Cache[T any] struct {
	Name string
}

func (c *Cache[T]) Serialize(target T) error {
	errCode := apperr.CacheSerialize

	_ = os.Mkdir(CacheDir, 0o755)

	f, err := os.Create(filepath.Join(CacheDir, c.Name+".bin"))
	if err != nil {
		return failure.Translate(err, errCode)
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	if err := enc.Encode(target); err != nil {
		return failure.Translate(err, errCode)
	}

	return nil
}

func (c *Cache[T]) Deserialize() (T, error) {
	errCode := apperr.CacheDeserialize

	_ = os.Mkdir(CacheDir, 0o755)

	var object T
	f, err := os.Open(filepath.Join(CacheDir, c.Name+".bin"))
	if err != nil {
		return object, failure.Translate(err, errCode)
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	if err := dec.Decode(&object); err != nil {
		return object, failure.Translate(err, errCode)
	}

	return object, nil
}
