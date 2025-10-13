package infra

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/morikuni/failure"
)

func load[T any](path string) (*T, error) {
	errCtx := failure.Context{"path": path}

	f, err := os.ReadFile(path)
	if err != nil {
		return nil, failure.Wrap(err, errCtx)
	}
	errCtx["data"] = string(f)

	var result T
	if err = json.Unmarshal(f, &result); err != nil {
		return nil, failure.Wrap(err, errCtx)
	}

	return &result, nil
}

func save[T any](path string, data T) error {
	//nolint:errchkjson
	b, _ := json.Marshal(data)
	errCtx := failure.Context{"path": path, "data": string(b)}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return failure.Wrap(err, errCtx)
	}

	file, err := os.Create(path)
	if err != nil {
		return failure.Wrap(err, errCtx)
	}
	//nolint:errcheck
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", strings.Repeat(" ", 2))

	if err = encoder.Encode(data); err != nil {
		return failure.Wrap(err, errCtx)
	}

	return nil
}
