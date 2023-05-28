package service

import (
	"changeme/backend/apperr"
	"changeme/backend/infra"
	"changeme/backend/vo"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type Config struct {
	configRepo    infra.ConfigInterface
	wargamingRepo infra.WargamingInterface
}

func NewConfig(
	configRepo infra.ConfigInterface,
	wargamingRepo infra.WargamingInterface,
) *Config {
	return &Config{
		configRepo:    configRepo,
		wargamingRepo: wargamingRepo,
	}
}

func (c *Config) User() (vo.UserConfig, error) {
	return c.configRepo.User()
}

func (c *Config) UpdateUser(config vo.UserConfig) error {
	if err := c.validate(config); err != nil {
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

func (c *Config) validate(config vo.UserConfig) error {
	if _, err := os.Stat(filepath.Join(config.InstallPath, "WorldOfWarships.exe")); err != nil {
		return errors.WithStack(apperr.SrvCfg.InvalidInstallPath.WithRaw(apperr.ErrInvalidInstallPath))
	}

	c.wargamingRepo.SetAppID(config.Appid)
	if _, err := c.wargamingRepo.EncyclopediaInfo(); err != nil {
		return errors.WithStack(apperr.SrvCfg.InvalidAppID.WithRaw(apperr.ErrInvalidAppID))
	}

	// Same value as "font-size": https://developer.mozilla.org/ja/docs/Web/CSS/font-size
	if !lo.Contains([]string{"x-small", "small", "medium", "large", "x-large"}, config.FontSize) {
		return errors.WithStack(apperr.SrvCfg.InvalidFontSize.WithRaw(apperr.ErrInvalidFontSize))
	}

	return nil
}
