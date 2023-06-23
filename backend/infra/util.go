package infra

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func readJSON[T any](path string, defaultValue T) (T, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return defaultValue, err
	}

	if err := json.Unmarshal(f, &defaultValue); err != nil {
		return defaultValue, err
	}

	return defaultValue, nil
}

func writeJSON[T any](path string, target T) error {
	_ = os.MkdirAll(filepath.Dir(path), 0o755)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")

	return encoder.Encode(target)
}
