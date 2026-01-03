package infra

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"wfs/internal/data"

	"github.com/morikuni/failure"
)

type Persistence struct {
	basePath string
}

func NewPersistence(basePath string) *Persistence {
	return &Persistence{
		basePath: basePath,
	}
}

// IGN.
func (p *Persistence) LoadOwnIGN() (string, error) {
	return _readString(filepath.Join(p.basePath, "own_ign.txt"))
}
func (p *Persistence) SaveOwnIGN(ign string) error {
	return _writeString(filepath.Join(p.basePath, "own_ign.txt"), ign)
}

// Expected Stats.
func (p *Persistence) LoadExpectedStats() (data.ExpectedStats, error) {
	return _readJSON[data.ExpectedStats](filepath.Join(p.basePath, "expected_stats.json"))
}
func (p *Persistence) SaveExpectedStats(data data.ExpectedStats) error {
	return _writeJSON(filepath.Join(p.basePath, "expected_stats.json"), data)
}

// User Config.
func (p *Persistence) LoadUserConfig() (data.UserConfig, error) {
	return _readJSON[data.UserConfig](filepath.Join(p.basePath, "user_config.json"))
}
func (p *Persistence) SaveUserConfig(data data.UserConfig) error {
	return _writeJSON(filepath.Join(p.basePath, "user_config.json"), data)
}

// Private functions.

func _readString(path string) (string, error) {
	errCtx := failure.Context{"path": path}

	f, err := os.ReadFile(path)
	if err != nil {
		return "", failure.Wrap(err, errCtx)
	}

	errCtx["data"] = string(f)

	return string(f), nil
}

func _writeString(path string, data string) error {
	errCtx := failure.Context{"path": path, "data": data}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return failure.Wrap(err, errCtx)
	}

	if err := os.WriteFile(path, []byte(data), os.ModePerm); err != nil {
		return failure.Wrap(err, errCtx)
	}

	return nil
}

func _readJSON[T any](path string) (T, error) {
	errCtx := failure.Context{"path": path}

	var result T
	f, err := os.ReadFile(path)
	if err != nil {
		return result, failure.Wrap(err, errCtx)
	}
	errCtx["data"] = string(f)

	if err = json.Unmarshal(f, &result); err != nil {
		return result, failure.Wrap(err, errCtx)
	}

	return result, nil
}

func _writeJSON[T any](path string, data T) error {
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
