package service

import (
	"context"
	"os"
	"path/filepath"
	"wfs/backend/apperr"
	"wfs/backend/application/repository"
	"wfs/backend/application/vo"
	"wfs/backend/domain"

	"github.com/morikuni/failure"
	"github.com/skratchdot/open-golang/open"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const GameExeName = "WorldOfWarships.exe"

type Config struct {
	localFile           repository.LocalFileInterface
	wargaming           repository.WargamingInterface
	WindowGetSize       WindowGetSize
	WindowSetSize       WindowSetSize
	WindowGetPosition   WindowGetPosition
	WindowSetPosition   WindowSetPosition
	OpenDirectoryDialog OpenDirectoryDialog
	OpenWithDefaultApp  OpenWithDefaultApp
}

func NewConfig(
	localFile repository.LocalFileInterface,
	wargaming repository.WargamingInterface,
) *Config {
	return &Config{
		localFile:           localFile,
		wargaming:           wargaming,
		WindowGetSize:       runtime.WindowGetSize,
		WindowSetSize:       runtime.WindowSetSize,
		WindowGetPosition:   runtime.WindowGetPosition,
		WindowSetPosition:   runtime.WindowSetPosition,
		OpenDirectoryDialog: runtime.OpenDirectoryDialog,
		OpenWithDefaultApp:  open.Run,
	}
}

func (c *Config) User() (domain.UserConfig, error) {
	config, err := c.localFile.User()
	return config, failure.Wrap(err)
}

func (c *Config) ValidateRequired(
	installPath string,
	appid string,
) vo.RequiredConfigError {
	result := vo.RequiredConfigError{}

	if installPath == "" {
		result.InstallPath = apperr.EmptyInstallPath.ErrorCode()
	} else {
		if _, err := os.Stat(filepath.Join(installPath, GameExeName)); err != nil {
			result.InstallPath = apperr.InvalidInstallPath.ErrorCode()
		}
	}

	if appid == "" {
		result.AppID = apperr.EmptyAppID.ErrorCode()
	} else {
		if ok, _ := c.wargaming.Test(appid); !ok {
			result.AppID = apperr.InvalidAppID.ErrorCode()
		}
	}

	result.Valid = result.InstallPath == "" && result.AppID == ""
	return result
}

func (c *Config) UpdateRequired(
	installPath string,
	appid string,
) (vo.RequiredConfigError, error) {
	// validate
	validatedResult := c.ValidateRequired(installPath, appid)
	if !validatedResult.Valid {
		return validatedResult, nil
	}

	// Note: overwrite only required setting
	config, err := c.localFile.User()
	if err != nil {
		return validatedResult, failure.Wrap(err)
	}
	config.InstallPath = installPath
	config.Appid = appid

	// write
	err = c.localFile.UpdateUser(config)

	return validatedResult, failure.Wrap(err)
}

func (c *Config) UpdateOptional(config domain.UserConfig) error {
	// Note: exclulde required setting
	saved, err := c.localFile.User()
	if err != nil {
		return failure.Wrap(err)
	}
	config.InstallPath = saved.InstallPath
	config.Appid = saved.Appid

	// write
	err = c.localFile.UpdateUser(config)
	return failure.Wrap(err)
}

func (c *Config) ApplyAppConfig(appCtx context.Context) error {
	config, err := c.localFile.App()
	if err != nil {
		return failure.Wrap(err)
	}

	// Set window size and position
	window := config.Window
	if window.Width > 0 && window.Height > 0 {
		c.WindowSetSize(appCtx, window.Width, window.Height)
	}
	if window.X >= 0 || window.Y >= 0 {
		c.WindowSetPosition(appCtx, window.X, window.Y)
	}

	return nil
}

func (c *Config) SaveAppConfig(appCtx context.Context) error {
	config, _ := c.localFile.App()

	// Save windows size and position
	width, height := c.WindowGetSize(appCtx)
	config.Window.Width = width
	config.Window.Height = height
	x, y := c.WindowGetPosition(appCtx)
	config.Window.X = x
	config.Window.Y = y

	return failure.Wrap(c.localFile.UpdateApp(config))
}

func (c *Config) AlertPlayers() ([]domain.AlertPlayer, error) {
	players, err := c.localFile.AlertPlayers()
	return players, failure.Wrap(err)
}

func (c *Config) UpdateAlertPlayer(player domain.AlertPlayer) error {
	err := c.localFile.UpdateAlertPlayer(player)
	return failure.Wrap(err)
}

func (c *Config) RemoveAlertPlayer(accountID int) error {
	err := c.localFile.RemoveAlertPlayer(accountID)
	return failure.Wrap(err)
}

func (c *Config) SearchPlayer(prefix string) (domain.WGAccountList, error) {
	list, err := c.wargaming.AccountListForSearch(prefix)
	return list, failure.Wrap(err)
}

func (c *Config) SelectDirectory(appCtx context.Context) (string, error) {
	selected, err := c.OpenDirectoryDialog(appCtx, runtime.OpenDialogOptions{})
	if err != nil {
		return selected, failure.New(apperr.WailsError, failure.Messagef("%s", err.Error()))
	}

	return selected, nil
}

func (c *Config) OpenDirectory(path string) error {
	err := c.OpenWithDefaultApp(path)
	if err != nil {
		return failure.New(apperr.OpenDirectoryError, failure.Context{"path": path}, failure.Messagef("%s", err.Error()))
	}

	return nil
}
