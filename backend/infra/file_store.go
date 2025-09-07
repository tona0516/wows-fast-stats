package infra

import (
	"os"
	"path/filepath"
	"wfs/backend/data"

	"github.com/morikuni/failure"
)

type FileStore struct {
	basePath string
}

func NewFileStore(basePath string) *FileStore {
	return &FileStore{basePath: basePath}
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

func (fs FileStore) Files(childPath string) ([]string, error) {
	p := filepath.Join(fs.basePath, childPath)
	_ = os.MkdirAll(p, os.ModePerm)

	entries, err := os.ReadDir(p)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

func createFilePath(basePath string, filePath data.FileStorePath) string {
	return basePath + "/" + filePath.ToString()
}
