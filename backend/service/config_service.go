package service

import (
	"changeme/backend/repo"
	"changeme/backend/vo"
	"errors"
	"os"
	"path/filepath"
)

type ConfigService struct{}

func (c *ConfigService) ReadUserConfig() (vo.UserConfig, error) {
    var config vo.UserConfig
    configAdapter := repo.ConfigAdapter{}
    config, err := configAdapter.ReadUserConfig()
    return config, err
}

func (c *ConfigService) UpdateUserConfig(config vo.UserConfig) (vo.UserConfig, error) {
    configAdapter := repo.ConfigAdapter{}

    if err := validate(config); err != nil {
        return config, err
    }

    if err := configAdapter.UpdateUserConfig(config); err != nil {
        return config, err
    }

    return config, nil
}

func (c *ConfigService) ReadAppConfig() (vo.AppConfig, error) {
    var config vo.AppConfig
    configAdapter := repo.ConfigAdapter{}
    config, err := configAdapter.ReadAppConfig()
    return config, err
}

func (c *ConfigService) UpdateAppConfig(config vo.AppConfig) (vo.AppConfig, error) {
    configAdapter := repo.ConfigAdapter{}

    if err := configAdapter.UpdateAppConfig(config); err != nil {
        return config, err
    }

    return config, nil
}

func validate(config vo.UserConfig) error {
    if _, err := os.Stat(filepath.Join(config.InstallPath, "WorldOfWarships.exe")); err != nil {
        err := errors.New("選択したパスに「WorldOfWarships.exe」が存在しません。")
        return err
    }

    wargaming := repo.Wargaming{AppID: config.Appid}
    if _, err := wargaming.GetEncyclopediaInfo(); err != nil {
        err := errors.New("AppIDが間違っています。")
        return err
    }

    // Same value as "font-size": https://developer.mozilla.org/ja/docs/Web/CSS/font-size
    if !contains([]string{"x-small", "small", "medium", "large", "x-large"}, config.FontSize) {
        err := errors.New("不正な文字サイズです。")
        return err
    }

    return nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
