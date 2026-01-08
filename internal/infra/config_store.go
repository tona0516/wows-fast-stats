package infra

import (
	"path/filepath"
	"wfs/internal/data"
)

type ConfigStore struct {
	dir            string
	userConfigFile string
}

func NewConfigStore(dir string) *ConfigStore {
	return &ConfigStore{
		dir:            dir,
		userConfigFile: "user_config.json",
	}
}

func (s *ConfigStore) UserConfig() (data.UserConfig, error) {
	return readJSON[data.UserConfig](filepath.Join(s.dir, s.userConfigFile))
}

func (s *ConfigStore) SetUserConfig(data data.UserConfig) error {
	return writeJSON(filepath.Join(s.dir, s.userConfigFile), data)
}
