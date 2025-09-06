package infra

import (
	"os"
	"wfs/backend/data"

	"github.com/morikuni/failure"
)

type FileStore struct {
	basePath string
}

func NewFileStore(path string) *FileStore {
	return &FileStore{basePath: path}
}

func (fs FileStore) Put(path data.FileStorePath, value string) error {
	_ = os.MkdirAll(fs.basePath, os.ModePerm)

	file, err := os.OpenFile(createFilePath(fs.basePath, path),
		os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return failure.Wrap(err)
	}
	//nolint:errcheck
	defer file.Close()

	_, err = file.WriteString(value)

	return failure.Wrap(err)
}

func (fs FileStore) Get(path data.FileStorePath) (string, error) {
	p := createFilePath(fs.basePath, path)

	data, err := os.ReadFile(p)
	if err != nil {
		return "", failure.Wrap(err)
	}

	return string(data), nil
}

func (fs FileStore) Delete(path data.FileStorePath) error {
	return os.Remove(createFilePath(fs.basePath, path))
}

func (fs FileStore) Keys() ([]string, error) {
	files, err := os.ReadDir(fs.basePath)
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

func createFilePath(basePath string, filePath data.FileStorePath) string {
	return basePath + "/" + filePath.ToString()
}
