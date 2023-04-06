package repo

import (
	"encoding/gob"
	"os"
)

type Cache[T any] struct {
	FileName string
}

func (s *Cache[T]) Serialize(object T) error {
	os.Mkdir("cache", 0755)

	f, err := os.Create("cache/" + s.FileName)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)

	if err := enc.Encode(object); err != nil {
		return err
	}

	return nil
}

func (s *Cache[T]) Deserialize() (T, error) {
	var object T

	f, err := os.Open("cache/" + s.FileName)
	if err != nil {
		return object, err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	if err := dec.Decode(&object); err != nil {
		return object, err
	}

	return object, nil
}
