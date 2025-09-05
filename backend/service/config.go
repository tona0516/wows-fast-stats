package service

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	fileStore           repository.FileStoreInterface
	logger              repository.LoggerInterface
	OpenDirectoryDialog openDirectoryDialogFunc
	OpenWithDefaultApp  openWithDefaultAppFunc
}

func NewConfig(
	localFile repository.LocalFileInterface,
	wargaming repository.WargamingInterface,
	fileStore repository.FileStoreInterface,
	logger repository.LoggerInterface,
) *Config {
	return &Config{
		localFile:           localFile,
		wargaming:           wargaming,
		fileStore:           fileStore,
		logger:              logger,
		OpenDirectoryDialog: runtime.OpenDirectoryDialog,
		OpenWithDefaultApp: func(input string) error {
			return exec.Command("explorer", input).Start()
		},
	}
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

func (c *Config) InstallPath() string {
	installPath, err := c.fileStore.Get(data.InstallPathKey)
	if err != nil {
		return ""
	}

	return installPath
}

func (c *Config) UpdateInstallPath(path string) error {
	// validate
	if err := c.ValidateInstallPath(path); err != nil {
		return err
	}

	return c.fileStore.Put(data.InstallPathKey, path)
}

func (c *Config) SendReport() bool {
	value, err := c.fileStore.Get(data.SendReportKey)
	if err != nil {
		return false
	}

	if value == "" {
		return false
	}

	return value == "1"
}

func (c *Config) UpdateSendReport(sendReport bool) error {
	var value string
	if sendReport {
		value = "1"
	} else {
		value = "0"
	}

	return c.fileStore.Put(data.SendReportKey, value)
}

func (c *Config) AlertPlayers() ([]data.AlertPlayer, error) {
	players := make([]data.AlertPlayer, 0)

	keys, err := c.fileStore.Keys()
	if err != nil {
		return nil, failure.Wrap(err)
	}

	for _, v := range keys {
		if !strings.HasPrefix(v, data.AlertPlayerKeyPrefix.ToString()) {
			continue
		}

		playerBytes, err := c.fileStore.Get(data.FileStoreKey(v))
		if err != nil {
			continue
		}

		var player data.AlertPlayer
		if err := json.Unmarshal([]byte(playerBytes), &player); err != nil {
			continue
		}

		players = append(players, player)
	}

	return players, nil
}

func (c *Config) UpdateAlertPlayer(player data.AlertPlayer) error {
	playerBytes, err := json.Marshal(player)
	if err != nil {
		return failure.Wrap(err)
	}

	key := data.AlertPlayerKeyPrefix.ToAlertPlayerKey(player.AccountID)
	if err := c.fileStore.Put(key, string(playerBytes)); err != nil {
		return failure.Wrap(err)
	}

	return nil
}

func (c *Config) RemoveAlertPlayer(accountID int) error {
	if err := c.fileStore.Delete(data.AlertPlayerKeyPrefix.ToAlertPlayerKey(accountID)); err != nil {
		return failure.Wrap(err)
	}

	return nil
}

func (c *Config) SearchPlayer(prefix string) (data.WGAccountList, error) {
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
