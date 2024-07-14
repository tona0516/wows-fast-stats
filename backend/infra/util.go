package infra

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func readJSON[T any](path string) (T, error) {
	var result T
	b, err := readFile(path)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(b, &result)
	return result, err
}

func writeFile(path string, target []byte) error {
	return os.WriteFile(path, target, 0o644)
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
	err = encoder.Encode(target)
	return err
}
