package infra

import (
	"path/filepath"
	"wfs/internal/data"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type ConfigStore interface {
	UserConfig() (data.UserConfig, error)
	SetUserConfig(data data.UserConfig) error
}

type configStore struct {
	dir            string
	userConfigFile string
}

func NewConfigStore(dir string) ConfigStore {
	return &configStore{
		dir:            dir,
		userConfigFile: "user_config.json",
	}
}

func (s *configStore) UserConfig() (data.UserConfig, error) {
	return readJSON[data.UserConfig](filepath.Join(s.dir, s.userConfigFile))
}

func (s *configStore) SetUserConfig(data data.UserConfig) error {
	return writeJSON(filepath.Join(s.dir, s.userConfigFile), data)
}
