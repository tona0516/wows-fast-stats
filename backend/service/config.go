package service

import (
	"changeme/backend/apperr"
	"changeme/backend/infra"
	"changeme/backend/vo"
	"os"
	"path/filepath"

	"github.com/morikuni/failure"
	"golang.org/x/exp/slices"
)

type Config struct {
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
	return c.configRepo.App()
}

func (c *Config) UpdateApp(config vo.AppConfig) error {
	return c.configRepo.UpdateApp(config)
}

func validate(config vo.UserConfig) error {
	if _, err := os.Stat(filepath.Join(config.InstallPath, "WorldOfWarships.exe")); err != nil {
		return failure.New(
			apperr.CfgSvInvalidInstallPath,
			failure.Message("選択したフォルダに「WorldOfWarships.exe」が存在しません。"),
		)
	}

	wargaming := infra.Wargaming{AppID: config.Appid}
	if _, err := wargaming.EncyclopediaInfo(); err != nil {
		return failure.New(
			apperr.CfgSvInvalidAppID,
			failure.Message("WG APIと通信できません。AppIDが間違っている可能性があります。"),
		)
	}

	// Same value as "font-size": https://developer.mozilla.org/ja/docs/Web/CSS/font-size
	if !slices.Contains([]string{"x-small", "small", "medium", "large", "x-large"}, config.FontSize) {
		return failure.New(
			apperr.CfgSvInvalidFontSize,
			failure.Message("不正な文字サイズです。"),
		)
	}

	return nil
}
