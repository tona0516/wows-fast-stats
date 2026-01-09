package infra

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/morikuni/failure"
)

const (
	pathKey = "path"
	dataKey = "data"
)

func readString(path string) (string, error) {
	errCtx := failure.Context{pathKey: path}

	f, err := os.ReadFile(path)
	if err != nil {
		return "", failure.Wrap(err, errCtx)
	}

	return string(f), nil
}

func writeString(path string, data string) error {
	errCtx := failure.Context{pathKey: path, dataKey: data}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return failure.Wrap(err, errCtx)
	}

	if err := os.WriteFile(path, []byte(data), os.ModePerm); err != nil {
		return failure.Wrap(err, errCtx)
	}

	return nil
}

func readJSON[T any](path string) (T, error) {
	errCtx := failure.Context{pathKey: path}

	var result T
	f, err := os.ReadFile(path)
	if err != nil {
		return result, failure.Wrap(err, errCtx)
	}
	errCtx[dataKey] = string(f)

	if err = json.Unmarshal(f, &result); err != nil {
		return result, failure.Wrap(err, errCtx)
	}

	return result, nil
}

func writeJSON[T any](path string, data T) error {
	//nolint:errchkjson
	b, _ := json.Marshal(data)
	errCtx := failure.Context{pathKey: path, dataKey: string(b)}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
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
