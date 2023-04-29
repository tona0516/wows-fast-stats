package service

import (
	"changeme/backend/infra"
	"changeme/backend/vo"
	"errors"
	"os"
	"path/filepath"
)

type Config struct{}

func (c *Config) User() (vo.UserConfig, error) {
    var config vo.UserConfig
    configRepo := infra.Config{}
    config, err := configRepo.User()
    return config, err
}

func (c *Config) UpdateUser(config vo.UserConfig) error {
    configRepo := infra.Config{}

    if err := validate(config); err != nil {
        return err
    }

    return configRepo.UpdateUser(config)
}

func (c *Config) App() (vo.AppConfig, error) {
    var config vo.AppConfig
    configRepo := infra.Config{}
    config, err := configRepo.App()
    return config, err
}

func (c *Config) UpdateApp(config vo.AppConfig) error {
    configRepo := infra.Config{}
    return configRepo.UpdateApp(config)
}

func validate(config vo.UserConfig) error {
    if _, err := os.Stat(filepath.Join(config.InstallPath, "WorldOfWarships.exe")); err != nil {
        err := errors.New("選択したフォルダに「WorldOfWarships.exe」が存在しません。")
        return err
    }

    wargaming := infra.Wargaming{AppID: config.Appid}
    if _, err := wargaming.EncyclopediaInfo(); err != nil {
        err := errors.New("WG APIと通信できません。AppIDが間違っている可能性があります。")
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
