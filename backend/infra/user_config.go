package infra

import (
	"wfs/backend/domain"
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
	return load[domain.UserConfig](u.basePath)
}

func (u *UserConfig) Save(data domain.UserConfig) error {
	return save(u.basePath, data)
}
