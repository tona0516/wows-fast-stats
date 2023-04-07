package service

import (
	"changeme/backend/repo"
	"changeme/backend/vo"
)

type ConfigService struct {
}

func (c *ConfigService) Read() (vo.Config, error) {
    var config vo.Config
    configAdapter := repo.ConfigAdapter{}
    config, err := configAdapter.Read()
    return config, err
}

func (c *ConfigService) Update(config vo.Config) (vo.Config, error) {
    // TODO validation
    configAdapter := repo.ConfigAdapter{}
    if err := configAdapter.Update(config); err != nil {
        return config, err
    }
    return config, nil
}
