package usecase

import (
	"context"
	"os"
	"path/filepath"
	"wfs/backend/apperr"
	"wfs/backend/domain/model"
	"wfs/backend/domain/repository"

	"github.com/morikuni/failure"
	"github.com/skratchdot/open-golang/open"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const GameExeName = "WorldOfWarships.exe"

type Config struct {
	localFile           repository.LocalFileInterface
	wargaming           repository.WargamingInterface
	storage             repository.StorageInterface
	logger              repository.LoggerInterface
	OpenDirectoryDialog openDirectoryDialogFunc
	OpenWithDefaultApp  openWithDefaultAppFunc
}

func NewConfig(
	localFile repository.LocalFileInterface,
	wargaming repository.WargamingInterface,
	storage repository.StorageInterface,
	logger repository.LoggerInterface,
) *Config {
	return &Config{
		localFile:           localFile,
		wargaming:           wargaming,
		storage:             storage,
		logger:              logger,
		OpenDirectoryDialog: runtime.OpenDirectoryDialog,
		OpenWithDefaultApp:  open.Run,
	}
}

func (c *Config) User() (model.UserConfigV2, error) {
	return c.storage.UserConfigV2()
}

func (c *Config) ValidateRequired(
	installPath string,
	appid string,
) model.RequiredConfigError {
	result := model.RequiredConfigError{}

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
) (model.RequiredConfigError, error) {
	// validate
	validatedResult := c.ValidateRequired(installPath, appid)
	if !validatedResult.Valid {
		return validatedResult, nil
	}

	// Note: overwrite only required setting
	config, err := c.storage.UserConfigV2()
	if err != nil {
		return validatedResult, err
	}
	config.InstallPath = installPath
	config.Appid = appid

	// write
	err = c.storage.WriteUserConfigV2(config)

	return validatedResult, err
}

func (c *Config) UpdateOptional(config model.UserConfigV2) error {
	// Note: exclulde required setting
	saved, err := c.storage.UserConfigV2()
	if err != nil {
		return err
	}
	config.InstallPath = saved.InstallPath
	config.Appid = saved.Appid

	// write
	err = c.storage.WriteUserConfigV2(config)
	return err
}

func (c *Config) AlertPlayers() ([]model.AlertPlayer, error) {
	players, err := c.storage.AlertPlayers()
	return players, err
}

func (c *Config) UpdateAlertPlayer(player model.AlertPlayer) error {
	players, err := c.storage.AlertPlayers()
	if err != nil {
		return err
	}

	var isMatched bool
	for i, v := range players {
		if player.AccountID == v.AccountID {
			players[i] = player
			isMatched = true
			break
		}
	}

	if !isMatched {
		players = append(players, player)
	}

	return c.storage.WriteAlertPlayers(players)
}

func (c *Config) RemoveAlertPlayer(accountID int) error {
	players, err := c.storage.AlertPlayers()
	if err != nil {
		return err
	}

	var isMatched bool
	for i, v := range players {
		if accountID == v.AccountID {
			players = players[:i+copy(players[i:], players[i+1:])]
			isMatched = true
			break
		}
	}

	if !isMatched {
		return nil
	}

	return c.storage.WriteAlertPlayers(players)
}

func (c *Config) SearchPlayer(prefix string) (model.WGAccountList, error) {
	return c.wargaming.AccountListForSearch(prefix)
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
