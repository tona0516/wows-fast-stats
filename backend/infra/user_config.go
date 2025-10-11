package infra

import (
	"encoding/json"
	"os"
	"path/filepath"
	"wfs/backend/domain"

	"github.com/morikuni/failure"
)

type UserConfig struct {
	basePath string
}

func NewUserConfig(basePath string) *UserConfig {
	return &UserConfig{
		basePath: basePath,
	}
}

func (u *UserConfig) Load() (*domain.UserConfig, error) {
	errCtx := failure.Context{"path": u.basePath}

	f, err := os.ReadFile(u.basePath)
	if err != nil {
		return nil, failure.Wrap(err, errCtx)
	}
	errCtx["data"] = string(f)

	var result domain.UserConfig
	if err = json.Unmarshal(f, &result); err != nil {
		return nil, failure.Wrap(err, errCtx)
	}

	return &result, nil
}

func (u *UserConfig) Save(data domain.UserConfig) error {
	//nolint:errchkjson
	b, _ := json.Marshal(data)
	errCtx := failure.Context{"path": u.basePath, "data": string(b)}

	_ = os.MkdirAll(filepath.Dir(u.basePath), 0o755)
	f, err := os.Create(u.basePath)
	if err != nil {
		return failure.Wrap(err, errCtx)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(data); err != nil {
		return failure.Wrap(err, errCtx)
	}

	return nil
}
