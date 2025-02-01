package service

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"wfs/backend/apperr"
	"wfs/backend/domain/model"
	"wfs/backend/domain/repository"

	"github.com/morikuni/failure"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const GameExeName = "WorldOfWarships.exe"

type Config struct {
	accountFetcher      repository.AccountFetcher
	userConfig          repository.UserConfigStore
	alertPlayer         repository.AlertPlayerStore
	OpenDirectoryDialog openDirectoryDialogFunc
	OpenWithDefaultApp  openWithDefaultAppFunc
}

func NewConfig(
	accountFetcher repository.AccountFetcher,
	userConfig repository.UserConfigStore,
	alertPlayer repository.AlertPlayerStore,
) *Config {
	return &Config{
		accountFetcher:      accountFetcher,
		userConfig:          userConfig,
		alertPlayer:         alertPlayer,
		OpenDirectoryDialog: runtime.OpenDirectoryDialog,
		OpenWithDefaultApp: func(input string) error {
			return exec.Command("explorer", input).Start()
		},
	}
}

func (c *Config) User() (model.UserConfigV2, error) {
	return c.userConfig.GetV2()
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

func (c *Config) UpdateInstallPath(path string) (model.UserConfigV2, error) {
	var config model.UserConfigV2

	// validate
	if err := c.ValidateInstallPath(path); err != nil {
		return config, err
	}

	// Note: overwrite only required setting
	config, err := c.userConfig.GetV2()
	if err != nil {
		return config, err
	}
	config.InstallPath = path

	// write
	return config, c.userConfig.SaveV2(config)
}

func (c *Config) UpdateOptional(config model.UserConfigV2) error {
	// Note: exclulde required setting
	saved, err := c.userConfig.GetV2()
	if err != nil {
		return err
	}
	config.InstallPath = saved.InstallPath

	// write
	err = c.userConfig.SaveV2(config)
	return err
}

func (c *Config) AlertPlayers() ([]model.AlertPlayer, error) {
	players, err := c.alertPlayer.GetV1()
	return players, err
}

func (c *Config) UpdateAlertPlayer(player model.AlertPlayer) ([]model.AlertPlayer, error) {
	var players []model.AlertPlayer

	players, err := c.alertPlayer.GetV1()
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

	return players, c.alertPlayer.SaveV1(players)
}

func (c *Config) RemoveAlertPlayer(accountID int) ([]model.AlertPlayer, error) {
	var players []model.AlertPlayer

	players, err := c.alertPlayer.GetV1()
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

	return players, c.alertPlayer.SaveV1(players)
}

func (c *Config) SearchPlayer(prefix string) model.Accounts {
	result, _ := c.accountFetcher.Search(prefix)
	return result
}

func (c *Config) SelectDirectory(appCtx context.Context) (string, error) {
	selected, err := c.OpenDirectoryDialog(appCtx, runtime.OpenDialogOptions{})
	if err != nil {
		return selected, failure.Translate(err, apperr.WailsError)
	}

	return selected, nil
}

func (c *Config) OpenDirectory(path string) error {
	err := c.OpenWithDefaultApp(path)
	if err != nil {
		return failure.Translate(err, apperr.OpenDirectoryError, failure.Context{"path": path})
	}

	return nil
}
