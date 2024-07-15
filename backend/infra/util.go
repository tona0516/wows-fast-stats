package infra

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/morikuni/failure"
)

func readFile(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	return b, failure.Wrap(err)
}

func readJSON[T any](path string) (T, error) {
	var result T

	b, err := readFile(path)
	if err != nil {
		return result, failure.Wrap(err)
	}

	err = json.Unmarshal(b, &result)
	return result, failure.Wrap(err)
}

func writeFile(path string, target []byte) error {
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	//nolint:gosec
	err := os.WriteFile(path, target, 0o644)
	return failure.Wrap(err)
}

func writeJSON[T any](path string, target T) error {
	b, err := json.Marshal(target)
	if err != nil {
		return failure.Wrap(err)
	}

	return writeFile(path, b)
}
