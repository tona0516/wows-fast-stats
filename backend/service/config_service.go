package service

import (
	"changeme/backend/repo"
	"changeme/backend/vo"
)

type ConfigService struct {
}

func (c *ConfigService) Read() (vo.UserConfig, error) {
    var config vo.UserConfig
    configAdapter := repo.ConfigAdapter{}
    config, err := configAdapter.Read()
    return config, err
}

func (c *ConfigService) Update(config vo.UserConfig) (vo.UserConfig, error) {
    // TODO validation
    configAdapter := repo.ConfigAdapter{}
    if err := configAdapter.Update(config); err != nil {
        return config, err
    }
    return config, nil
}
