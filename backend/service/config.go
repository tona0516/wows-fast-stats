package service

import (
	"changeme/backend/infra"
	"changeme/backend/vo"
	"os"
	"path/filepath"

	"golang.org/x/exp/slices"
    "github.com/pkg/errors"
)

type Config struct{
    configRepo infra.Config
}

func NewConfig(configRepo infra.Config) *Config {
    return &Config{
        configRepo: configRepo,
    }
}

func (c *Config) User() (vo.UserConfig, error) {
    return c.configRepo.User()
}

func (c *Config) UpdateUser(config vo.UserConfig) error {
    if err := validate(config); err != nil {
        return err
    }

    return c.configRepo.UpdateUser(config)
}

func (c *Config) App() (vo.AppConfig, error) {
    return  c.configRepo.App()
}

func (c *Config) UpdateApp(config vo.AppConfig) error {
    return c.configRepo.UpdateApp(config)
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
    if !slices.Contains([]string{"x-small", "small", "medium", "large", "x-large"}, config.FontSize) {
        err := errors.New("不正な文字サイズです。")
        return err
    }

    return nil
}
