package infra

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"wfs/backend/core"

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
		return "", failure.Translate(err, core.ErrStringRead, errCtx)
	}

	return string(f), nil
}

func writeString(path string, value string) error {
	errCtx := failure.Context{pathKey: path, dataKey: value}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return failure.Translate(err, core.ErrStringWrite, errCtx)
	}

	if err := os.WriteFile(path, []byte(value), os.ModePerm); err != nil {
		return failure.Translate(err, core.ErrStringWrite, errCtx)
	}

	return nil
}

func readJSON[T any](path string) (T, error) {
	errCtx := failure.Context{pathKey: path}

	var result T
	f, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return result, failure.Translate(err, core.ErrJSONNotFound, errCtx)
		}
		return result, failure.Translate(err, core.ErrJSONRead, errCtx)
	}
	errCtx[dataKey] = string(f)

	if err = json.Unmarshal(f, &result); err != nil {
		return result, failure.Translate(err, core.ErrJSONRead, errCtx)
	}

	return result, nil
}

func writeJSON[T any](path string, value T) error {
	//nolint:errchkjson
	b, _ := json.Marshal(value)
	errCtx := failure.Context{pathKey: path, dataKey: string(b)}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return failure.Translate(err, core.ErrJSONWrite, errCtx)
	}

	file, err := os.Create(path)
	if err != nil {
		return failure.Translate(err, core.ErrJSONWrite, errCtx)
	}
	//nolint:errcheck
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", strings.Repeat(" ", 2))

	if err = encoder.Encode(value); err != nil {
		return failure.Translate(err, core.ErrJSONWrite, errCtx)
	}

	return nil
}
