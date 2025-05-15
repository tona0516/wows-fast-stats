package infra

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"wfs/backend/apperr"

	"github.com/morikuni/failure"
)

func readJSON[T any](path string, defaulValue T) (T, error) {
	errCtx := failure.Context{"path": path}

	f, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return defaulValue, failure.New(apperr.FileNotExist, errCtx)
		}
		return defaulValue, failure.Wrap(err, errCtx)
	}
	errCtx["target"] = string(f)

	err = json.Unmarshal(f, &defaulValue)
	return defaulValue, failure.Wrap(err, errCtx)
}

func writeJSON[T any](path string, target T) error {
	//nolint:errchkjson
	marshaled, _ := json.Marshal(target)
	errCtx := failure.Context{"path": path, "target": string(marshaled)}

	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	f, err := os.Create(path)
	if err != nil {
		return failure.Wrap(err, errCtx)
	}
	//nolint:errcheck
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(target)
	return failure.Wrap(err, errCtx)
}
