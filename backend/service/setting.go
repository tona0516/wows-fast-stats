package service

import (
	"context"
	"encoding/json"
	"errors"
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

type Setting struct {
	localFile           repository.LocalFileInterface
	wargaming           repository.WargamingInterface
	fileStore           repository.FileStoreInterface
	logger              repository.LoggerInterface
	OpenDirectoryDialog openDirectoryDialogFunc
	OpenWithDefaultApp  openWithDefaultAppFunc
}

func NewSetting(
	localFile repository.LocalFileInterface,
	wargaming repository.WargamingInterface,
	fileStore repository.FileStoreInterface,
	logger repository.LoggerInterface,
) *Setting {
	return &Setting{
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

func (s *Setting) RequiredSetting() (data.RequiredSetting, error) {
	text, err := s.fileStore.Get(data.FileNameRequiredSetting)
	if errors.Is(err, os.ErrNotExist) {
		return data.RequiredSetting{Version: 1, InstallPath: ""}, nil
	}

	if err != nil {
		return data.RequiredSetting{}, failure.Wrap(err)
	}

	var setting data.RequiredSetting
	if err := json.Unmarshal([]byte(text), &setting); err != nil {
		return data.RequiredSetting{}, failure.Wrap(err)
	}

	return setting, nil
}

func (s *Setting) UpdateRequiredSetting(setting data.RequiredSetting) error {
	settingBytes, err := json.Marshal(setting)
	if err != nil {
		return failure.Wrap(err)
	}

	if err := s.fileStore.Put(data.FileNameRequiredSetting, string(settingBytes)); err != nil {
		return failure.Wrap(err)
	}

	return nil
}

func (s *Setting) ValidateInstallPath(path string) error {
	if path == "" {
		return failure.New(apperr.EmptyInstallPath)
	}

	if _, err := os.Stat(filepath.Join(path, GameExeName)); err != nil {
		return failure.New(apperr.InvalidInstallPath)
	}

	return nil
}

func (s *Setting) OptionalSetting() (data.OptionalSetting, error) {
	text, err := s.fileStore.Get(data.FileNameOptionalSetting)
	if errors.Is(err, os.ErrNotExist) {
		return data.OptionalSetting{Version: 1, ZoomRate: 100, StatsExtra: "pvp_all", IsSendReport: true}, nil
	}

	if err != nil {
		return data.OptionalSetting{}, failure.Wrap(err)
	}

	var setting data.OptionalSetting
	if err := json.Unmarshal([]byte(text), &setting); err != nil {
		return data.OptionalSetting{}, failure.Wrap(err)
	}

	return setting, nil
}

func (s *Setting) UpdateOptionalSetting(setting data.OptionalSetting) error {
	settingBytes, err := json.Marshal(setting)
	if err != nil {
		return failure.Wrap(err)
	}

	if err := s.fileStore.Put(data.FileNameOptionalSetting, string(settingBytes)); err != nil {
		return failure.Wrap(err)
	}

	return nil
}

func (s *Setting) BasicColumnSetting() (data.BasicColumnSetting, error) {
	text, err := s.fileStore.Get(data.FileNameBasicColumnSetting)
	if errors.Is(err, os.ErrNotExist) {
		return data.BasicColumnSetting{
			Version: 1,
			Ship: data.ShipColumnSetting{
				EnableNationFlag: true,
				IsColored:        false,
			},
			Player: data.PlayerColumnSetting{
				EnableNationFlag: false,
				ColorPattern:     "none",
			},
		}, nil
	}

	if err != nil {
		return data.BasicColumnSetting{}, failure.Wrap(err)
	}

	var setting data.BasicColumnSetting
	if err := json.Unmarshal([]byte(text), &setting); err != nil {
		return data.BasicColumnSetting{}, failure.Wrap(err)
	}

	return setting, nil
}

func (s *Setting) UpdateBasicColumnSetting(setting data.BasicColumnSetting) error {
	settingBytes, err := json.Marshal(setting)
	if err != nil {
		return failure.Wrap(err)
	}

	if err := s.fileStore.Put(data.FileNameBasicColumnSetting, string(settingBytes)); err != nil {
		return failure.Wrap(err)
	}

	return nil
}

func (s *Setting) StatsColumnSettings() (data.StatsColumnSettings, error) {
	text, err := s.fileStore.Get(data.FileNameStatsColumnSetting)
	if errors.Is(err, os.ErrNotExist) {
		return data.StatsColumnSettings{
			Version: 1,
			Battles: data.StatsColumnSetting{
				IsShowShip:    true,
				IsShowOverall: true,
				Digit:         0,
			},
			Damage: data.StatsColumnSetting{
				IsShowShip:    true,
				IsShowOverall: true,
				Digit:         0,
			},
			MaxDamage: data.StatsColumnSetting{
				IsShowShip:    true,
				IsShowOverall: true,
				Digit:         0,
			},
			WinRate: data.StatsColumnSetting{
				IsShowShip:    true,
				IsShowOverall: true,
				Digit:         1,
			},
			SurvivedRate: data.StatsColumnSetting{
				IsShowShip:    false,
				IsShowOverall: false,
				Digit:         1,
			},
			KdRate: data.StatsColumnSetting{
				IsShowShip:    false,
				IsShowOverall: false,
				Digit:         2,
			},
			Kill: data.StatsColumnSetting{
				IsShowShip:    true,
				IsShowOverall: true,
				Digit:         2,
			},
			Exp: data.StatsColumnSetting{
				IsShowShip:    true,
				IsShowOverall: true,
				Digit:         0,
			},
			PR: data.StatsColumnSetting{
				IsShowShip:    false,
				IsShowOverall: false,
				Digit:         0,
			},
			HitRate: data.StatsColumnSetting{
				IsShowShip:    false,
				IsShowOverall: false,
				Digit:         1,
			},
			PlanesKilled: data.StatsColumnSetting{
				IsShowShip:    false,
				IsShowOverall: false,
				Digit:         0,
			},
			PlatoonRate: data.StatsColumnSetting{
				IsShowShip:    false,
				IsShowOverall: false,
				Digit:         1,
			},
			EfficiencyBadge: data.StatsColumnSetting{
				IsShowShip:    false,
				IsShowOverall: false,
				Digit:         0,
			},
			ThreatLevel: data.StatsColumnSetting{
				IsShowShip:    false,
				IsShowOverall: false,
				Digit:         1,
			},
			AvgTier: data.StatsColumnSetting{
				IsShowShip:    false,
				IsShowOverall: false,
				Digit:         1,
			},
			UsingShipTypeRate: data.StatsColumnSetting{
				IsShowShip:    false,
				IsShowOverall: false,
				Digit:         1,
			},
			UsingTierRate: data.StatsColumnSetting{
				IsShowShip:    false,
				IsShowOverall: false,
				Digit:         1,
			},
		}, nil
	}

	if err != nil {
		return data.StatsColumnSettings{}, failure.Wrap(err)
	}

	var setting data.StatsColumnSettings
	if err := json.Unmarshal([]byte(text), &setting); err != nil {
		return data.StatsColumnSettings{}, failure.Wrap(err)
	}

	return setting, nil
}

func (s *Setting) UpdateStatsColumnSettings(setting data.StatsColumnSettings) error {
	settingBytes, err := json.Marshal(setting)
	if err != nil {
		return failure.Wrap(err)
	}

	if err := s.fileStore.Put(data.FileNameStatsColumnSetting, string(settingBytes)); err != nil {
		return failure.Wrap(err)
	}

	return nil
}

func (c *Setting) AlertPlayers() ([]data.AlertPlayer, error) {
	players := make([]data.AlertPlayer, 0)

	keys, err := c.fileStore.Files(string(data.FilePathAlertPlayer))
	if err != nil {
		return nil, failure.Wrap(err)
	}

	for _, v := range keys {
		playerBytes, err := c.fileStore.Get(data.FileStorePath(v))
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

func (c *Setting) UpdateAlertPlayer(player data.AlertPlayer) error {
	playerBytes, err := json.Marshal(player)
	if err != nil {
		return failure.Wrap(err)
	}

	filename := data.FilePathAlertPlayer.ToAlertPlayerFileName(player.AccountID)
	if err := c.fileStore.Put(filename, string(playerBytes)); err != nil {
		return failure.Wrap(err)
	}

	return nil
}

func (c *Setting) RemoveAlertPlayer(accountID int) error {
	filename := data.FilePathAlertPlayer.ToAlertPlayerFileName(accountID)

	if err := c.fileStore.Delete(filename); err != nil {
		return failure.Wrap(err)
	}

	return nil
}

func (c *Setting) SearchPlayer(prefix string) (data.WGAccountList, error) {
	return c.wargaming.AccountListForSearch(prefix)
}

func (c *Setting) SelectDirectory(appCtx context.Context) (string, error) {
	selected, err := c.OpenDirectoryDialog(appCtx, runtime.OpenDialogOptions{})
	if err != nil {
		return selected, failure.New(apperr.WailsError, failure.Messagef("%s", err.Error()))
	}

	return selected, nil
}

func (c *Setting) OpenDirectory(path string) error {
	err := c.OpenWithDefaultApp(path)
	if err != nil {
		return failure.New(apperr.OpenDirectoryError, failure.Context{"path": path}, failure.Messagef("%s", err.Error()))
	}

	return nil
}
