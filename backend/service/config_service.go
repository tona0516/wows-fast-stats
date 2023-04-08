package service

import (
	"changeme/backend/repo"
	"changeme/backend/vo"
	"errors"
	"os"
	"path/filepath"
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
    if _, err := os.Stat(filepath.Join(config.InstallPath, "WorldOfWarships.exe")); err != nil {
        err := errors.New("選択したパスに「WorldOfWarships.exe」が存在しません。")
        return config, err
    }

    wargaming := repo.Wargaming{AppID: config.Appid}
    if _, err := wargaming.GetEncyclopediaInfo(); err != nil {
        err := errors.New("AppIDが間違っています。")
        return config, err
    }

    configAdapter := repo.ConfigAdapter{}
    if err := configAdapter.Update(config); err != nil {
        return config, err
    }
    return config, nil
}
