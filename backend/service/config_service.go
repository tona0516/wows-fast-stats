package service

import (
	"changeme/backend/repo"
	"changeme/backend/vo"
	"errors"
	"os"
	"path/filepath"
)

type ConfigService struct{}

func (c *ConfigService) User() (vo.UserConfig, error) {
    var config vo.UserConfig
    configAdapter := repo.Config{}
    config, err := configAdapter.User()
    return config, err
}

func (c *ConfigService) UpdateUser(config vo.UserConfig) (vo.UserConfig, error) {
    configAdapter := repo.Config{}

    if err := validate(config); err != nil {
        return config, err
    }

    if err := configAdapter.UpdateUser(config); err != nil {
        return config, err
    }

    return config, nil
}

func (c *ConfigService) App() (vo.AppConfig, error) {
    var config vo.AppConfig
    configAdapter := repo.Config{}
    config, err := configAdapter.App()
    return config, err
}

func (c *ConfigService) UpdateApp(config vo.AppConfig) (vo.AppConfig, error) {
    configAdapter := repo.Config{}

    if err := configAdapter.UpdateApp(config); err != nil {
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
    if _, err := wargaming.EncyclopediaInfo(); err != nil {
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
