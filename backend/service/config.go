package service

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/repository"

	"github.com/morikuni/failure"
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
		OpenWithDefaultApp: func(input string) error {
			return exec.Command("explorer", input).Start()
		},
	}
}

func (c *Config) User() (data.UserConfigV2, error) {
	return c.storage.UserConfigV2()
}

func (c *Config) ValidateInstallPath(path string) error {
	if path == "" {
		return failure.New(apperr.EmptyInstallPath)
	}

	if _, err := os.Stat(filepath.Join(path, GameExeName)); err != nil {
		return failure.New(apperr.InvalidInstallPath)
	}

	return nil
}

func (c *Config) UpdateInstallPath(path string) (data.UserConfigV2, error) {
	var config data.UserConfigV2

	// validate
	if err := c.ValidateInstallPath(path); err != nil {
		return config, err
	}

	// Note: overwrite only required setting
	config, err := c.storage.UserConfigV2()
	if err != nil {
		return config, err
	}
	config.InstallPath = path

	// write
	return config, c.storage.WriteUserConfigV2(config)
}

func (c *Config) UpdateOptional(config data.UserConfigV2) error {
	// Note: exclulde required setting
	saved, err := c.storage.UserConfigV2()
	if err != nil {
		return err
	}
	config.InstallPath = saved.InstallPath

	// write
	err = c.storage.WriteUserConfigV2(config)
	return err
}

func (c *Config) AlertPlayers() ([]data.AlertPlayer, error) {
	players, err := c.storage.AlertPlayers()
	return players, err
}

func (c *Config) UpdateAlertPlayer(player data.AlertPlayer) ([]data.AlertPlayer, error) {
	var players []data.AlertPlayer

	players, err := c.storage.AlertPlayers()
	if err != nil {
		return players, err
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

	return players, c.storage.WriteAlertPlayers(players)
}

func (c *Config) RemoveAlertPlayer(accountID int) ([]data.AlertPlayer, error) {
	var players []data.AlertPlayer

	players, err := c.storage.AlertPlayers()
	if err != nil {
		return players, err
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
		return players, nil
	}

	return players, c.storage.WriteAlertPlayers(players)
}

func (c *Config) SearchPlayer(prefix string) data.WGAccountList {
	result, _ := c.wargaming.AccountListForSearch(prefix)
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
