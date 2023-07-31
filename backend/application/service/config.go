package service

import (
	"context"
	"os"
	"path/filepath"
	"wfs/backend/apperr"
	"wfs/backend/application/repository"
	"wfs/backend/application/vo"
	"wfs/backend/domain"

	"github.com/skratchdot/open-golang/open"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const GameExeName = "WorldOfWarships.exe"

type Config struct {
	localFile repository.LocalFileInterface
	wargaming repository.WargamingInterface
}

func NewConfig(
	localFile repository.LocalFileInterface,
	wargaming repository.WargamingInterface,
) *Config {
	return &Config{
		localFile: localFile,
		wargaming: wargaming,
	}
}

func (c *Config) User() (domain.UserConfig, error) {
	return c.localFile.User()
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

	// Note: overwrite only required setting
	config, err := c.localFile.User()
	if err != nil {
		return validatedResult, err
	}
	config.InstallPath = installPath
	config.Appid = appid

	// write
	return validatedResult, c.localFile.UpdateUser(config)
}

func (c *Config) UpdateOptional(config domain.UserConfig) error {
	// Note: exclulde required setting
	saved, err := c.localFile.User()
	if err != nil {
		return err
	}
	config.InstallPath = saved.InstallPath
	config.Appid = saved.Appid

	// write
	return c.localFile.UpdateUser(config)
}

func (c *Config) App() (vo.AppConfig, error) {
	return c.localFile.App()
}

func (c *Config) UpdateApp(config vo.AppConfig) error {
	return c.localFile.UpdateApp(config)
}

func (c *Config) AlertPlayers() ([]domain.AlertPlayer, error) {
	return c.localFile.AlertPlayers()
}

func (c *Config) UpdateAlertPlayer(player domain.AlertPlayer) error {
	return c.localFile.UpdateAlertPlayer(player)
}

func (c *Config) RemoveAlertPlayer(accountID int) error {
	return c.localFile.RemoveAlertPlayer(accountID)
}

func (c *Config) SearchPlayer(prefix string) (domain.WGAccountList, error) {
	return c.wargaming.AccountListForSearch(prefix)
}

func (c *Config) SelectDirectory(appCtx context.Context) (string, error) {
	selected, err := runtime.OpenDirectoryDialog(appCtx, runtime.OpenDialogOptions{})
	if err != nil {
		return "", apperr.New(apperr.ErrSelectDirectory, err)
	}

	return selected, nil
}

func (c *Config) OpenDirectory(path string) error {
	if err := open.Run(path); err != nil {
		return apperr.New(apperr.ErrOpenDirectory, err)
	}

	return nil
}

func (c *Config) validateRequired(
	installPath string,
	appid string,
) vo.ValidatedResult {
	result := vo.ValidatedResult{}

	if _, err := os.Stat(filepath.Join(installPath, GameExeName)); err != nil {
		result.InstallPath = apperr.New(apperr.ErrInvalidInstallPath, nil).Error()
	}

	if ok, err := c.wargaming.Test(appid); !ok {
		result.AppID = apperr.New(apperr.ErrInvalidAppID, err).Error()
	}

	return result
}
