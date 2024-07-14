package service

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/repository"

	"github.com/morikuni/failure"
	"github.com/skratchdot/open-golang/open"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const GameExeName = "WorldOfWarships.exe"

type Config struct {
	localFile           repository.LocalFileInterface
	wargaming           repository.WargamingInterface
	logger              repository.LoggerInterface
	OpenDirectoryDialog openDirectoryDialogFunc
	OpenWithDefaultApp  openWithDefaultAppFunc
}

func NewConfig(
	localFile repository.LocalFileInterface,
	wargaming repository.WargamingInterface,
	logger repository.LoggerInterface,
) *Config {
	return &Config{
		localFile:           localFile,
		wargaming:           wargaming,
		logger:              logger,
		OpenDirectoryDialog: runtime.OpenDirectoryDialog,
		OpenWithDefaultApp:  open.Run,
	}
}

func (c *Config) InitUserIfNeeded() {
	if _, err := c.localFile.UserConfigV2(); errors.Is(err, os.ErrNotExist) {
		_ = c.localFile.WriteUserConfigV2(data.DefaultUserConfigV2())
	}
}

func (c *Config) User() (data.UserConfigV2, error) {
	return c.localFile.UserConfigV2()
}

func (c *Config) ValidateRequired(
	installPath string,
	appid string,
) data.RequiredConfigError {
	result := data.RequiredConfigError{}

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
) (data.RequiredConfigError, error) {
	// validate
	validatedResult := c.ValidateRequired(installPath, appid)
	if !validatedResult.Valid {
		return validatedResult, nil
	}

	// Note: overwrite only required setting
	config, err := c.localFile.UserConfigV2()
	if err != nil {
		return validatedResult, err
	}
	config.InstallPath = installPath
	config.Appid = appid

	// write
	err = c.localFile.WriteUserConfigV2(config)

	return validatedResult, err
}

func (c *Config) UpdateOptional(config data.UserConfigV2) error {
	// Note: exclulde required setting
	saved, err := c.localFile.UserConfigV2()
	if err != nil {
		return err
	}
	config.InstallPath = saved.InstallPath
	config.Appid = saved.Appid

	// write
	err = c.localFile.WriteUserConfigV2(config)
	return err
}

func (c *Config) InitAlertPlayersIfNeeded() {
	if _, err := c.localFile.AlertPlayerV1(); errors.Is(err, os.ErrNotExist) {
		_ = c.localFile.WriteAlertPlayerV1([]data.AlertPlayer{})
	}
}

func (c *Config) AlertPlayers() ([]data.AlertPlayer, error) {
	players, err := c.localFile.AlertPlayerV1()
	return players, err
}

func (c *Config) UpdateAlertPlayer(player data.AlertPlayer) error {
	players, err := c.localFile.AlertPlayerV1()
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

	return c.localFile.WriteAlertPlayerV1(players)
}

func (c *Config) RemoveAlertPlayer(accountID int) error {
	players, err := c.localFile.AlertPlayerV1()
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

	return c.localFile.WriteAlertPlayerV1(players)
}

func (c *Config) SearchPlayer(prefix string) data.WGAccountList {
	config, err := c.localFile.UserConfigV2()
	if err != nil {
		return data.WGAccountList{}
	}

	appID := config.Appid
	if appID == "" {
		return data.WGAccountList{}
	}

	result, err := c.wargaming.AccountListForSearch(appID, prefix)
	if err != nil {
		return data.WGAccountList{}
	}

	return result
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
