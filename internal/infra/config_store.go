package infra

import (
	"path/filepath"
	"wfs/internal/config"
	"wfs/internal/data"

	"github.com/samber/do/v2"
)

type ConfigStore struct {
	dir            string
	userConfigFile string
}

func NewConfigStore(i do.Injector) (*ConfigStore, error) {
	config := do.MustInvoke[config.Config](i)
	return &ConfigStore{
		dir:            config.LocalFile.ConfigDir,
		userConfigFile: "user_config.json",
	}, nil
}

func (s *ConfigStore) UserConfig() (data.UserConfig, error) {
	return readJSON[data.UserConfig](filepath.Join(s.dir, s.userConfigFile))
}

func (s *ConfigStore) SetUserConfig(data data.UserConfig) error {
	return writeJSON(filepath.Join(s.dir, s.userConfigFile), data)
}
