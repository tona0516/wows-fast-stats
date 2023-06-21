package service

import (
	"changeme/backend/apperr"
	"changeme/backend/infra"
	"changeme/backend/vo"
	"os"
	"path/filepath"
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

func (c *Config) AlertPlayers() ([]vo.AlertPlayer, error) {
	return c.configRepo.AlertPlayers()
}

func (c *Config) UpdateAlertPlayer(player vo.AlertPlayer) error {
	return c.configRepo.UpdateAlertPlayer(player)
}

func (c *Config) RemoveAlertPlayer(accountID int) error {
	return c.configRepo.RemoveAlertPlayer(accountID)
}

func (c *Config) SearchPlayer(prefix string) (vo.WGAccountList, error) {
	return c.wargamingRepo.AccountListForSearch(prefix)
}

func (c *Config) validate(config vo.UserConfig) error {
	if _, err := os.Stat(filepath.Join(config.InstallPath, "WorldOfWarships.exe")); err != nil {
		return apperr.New(apperr.ValidateInvalidInstallPath, err)
	}

	c.wargamingRepo.SetAppID(config.Appid)
	if _, err := c.wargamingRepo.EncycInfo(); err != nil {
		return apperr.New(apperr.ValidateInvalidAppID, err)
	}

	// Same value as "font-size": https://developer.mozilla.org/ja/docs/Web/CSS/font-size
	var validFontSize bool
	for _, v := range vo.FontSizes {
		if v == config.FontSize {
			validFontSize = true
			break
		}
	}
	if !validFontSize {
		return apperr.New(apperr.ValidateInvalidFontSize, nil)
	}

	return nil
}
