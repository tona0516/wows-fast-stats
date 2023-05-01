package infra

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"strings"
)

const DIRECTORY = "cache"

type Cache[T any] struct {
	Prefix string
    GameVersion string
}

func (c *Cache[T]) Serialize(object T) error {
	_ = os.Mkdir(DIRECTORY, 0755)

    filename := c.Prefix + "_" + c.GameVersion + ".bin"
    path := filepath.Join(DIRECTORY, filename)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
    return enc.Encode(object)
}

func (c *Cache[T]) Deserialize() (T, error) {
	var object T

    filename := c.Prefix + "_" + c.GameVersion + ".bin"
    path := filepath.Join(DIRECTORY, filename)
	f, err := os.Open(path)
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

func (c *Cache[T]) RemoveOld() error {
    entries, err := os.ReadDir(DIRECTORY)
    if err != nil {
        return err
    }

    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }

        if !strings.HasPrefix(entry.Name(), c.Prefix) {
            continue
        }

        filename := c.Prefix + "_" + c.GameVersion + ".bin"
        if entry.Name() == filename {
            continue
        }

        os.Remove(filepath.Join(DIRECTORY, entry.Name()))
	}

    return nil
}
