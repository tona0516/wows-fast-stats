package service

import (
	"fmt"
	"os"
	"path/filepath"
	"wfs/backend/apperr"
	"wfs/backend/infra"
	"wfs/backend/vo"
)

const GameExeName = "WorldOfWarships.exe"

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

func (c *Config) UpdateRequired(
	installPath string,
	appid string,
) (vo.ValidatedResult, error) {
	// validate
	validatedResult := c.validateRequired(installPath, appid)
	if !validatedResult.Valid() {
		return validatedResult, nil
	}

	// Note: overerite only required setting
	config, err := c.configRepo.User()
	if err != nil {
		return validatedResult, err
	}
	config.InstallPath = installPath
	config.Appid = appid

	// write
	return validatedResult, c.configRepo.UpdateUser(config)
}

func (c *Config) UpdateOptional(config vo.UserConfig) error {
	// Note: exclulde required setting
	saved, err := c.configRepo.User()
	if err != nil {
		return err
	}
	config.InstallPath = saved.InstallPath
	config.Appid = saved.Appid

	// write
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

func (c *Config) validateRequired(
	installPath string,
	appid string,
) vo.ValidatedResult {
	result := vo.ValidatedResult{}

	if _, err := os.Stat(filepath.Join(installPath, GameExeName)); err != nil {
		result.InstallPath = apperr.ErrInvalidInstallPath.Error()
	}

	if ok, err := c.wargamingRepo.Test(appid); !ok {
		result.AppID = fmt.Sprintf("%s(%s)", apperr.ErrInvalidAppID, err.Error())
	}

	return result
}
