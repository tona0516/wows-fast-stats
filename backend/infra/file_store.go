package infra

import (
	"os"
	"wfs/backend/apperr"
	"wfs/backend/data"

	"github.com/morikuni/failure"
)

type FileStore struct {
	path string
}

func NewFileStore(path string) *FileStore {
	return &FileStore{path: path}
}

func (fs FileStore) Put(key data.FileStoreKey, value string) error {
	_ = os.MkdirAll(fs.path, 0o755)

	file, err := os.OpenFile(createFilePath(fs.path, key),
		os.O_WRONLY|os.O_CREATE, 0o755)
	if err != nil {
		return failure.Wrap(err)
	}
	//nolint:errcheck
	defer file.Close()

	_, err = file.WriteString(value)

	return failure.Wrap(err)
}

func (fs FileStore) Get(key data.FileStoreKey) (string, error) {
	p := createFilePath(fs.path, key)

	if _, err := os.Stat(p); os.IsNotExist(err) {
		return "", failure.New(apperr.FileNotExist)
	}

	data, err := os.ReadFile(p)
	if err != nil {
		return "", failure.Wrap(err)
	}

	return string(data), nil
}

func (fs FileStore) Remove(key data.FileStoreKey) error {
	return os.Remove(createFilePath(fs.path, key))
}

func (fs FileStore) Keys() ([]string, error) {
	files, err := os.ReadDir(fs.path)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	keys := make([]string, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			keys = append(keys, file.Name())
		}
	}
	return keys, nil
}

func createFilePath(basePath string, key data.FileStoreKey) string {
	return basePath + "/" + key.ToString()
}
