package service

import (
	"context"
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

type Config struct {
	persistence         repository.PersistenceInterface
	wargaming           repository.WargamingInterface
	logger              repository.LoggerInterface
	OpenDirectoryDialog openDirectoryDialogFunc
	OpenWithDefaultApp  openWithDefaultAppFunc
}

func NewSetting(
	persistence repository.PersistenceInterface,
	wargaming repository.WargamingInterface,
	logger repository.LoggerInterface,
) *Config {
	return &Config{
		persistence:         persistence,
		wargaming:           wargaming,
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

func (c *Config) TrySaveInstallPath(ctx context.Context) (bool, error) {
	selected, err := c.OpenDirectoryDialog(ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return false, failure.New(apperr.WailsError, failure.Messagef("%s", err.Error()))
	}

	if selected == "" {
		return false, nil
	}

	cfg, err := c.GetUserConfig()
	if err != nil {
		return false, failure.Wrap(err)
	}

	cfg.InstallPath = selected

	if err = c.SaveUserConfig(cfg); err != nil {
		return false, failure.Wrap(err)
	}

	return true, nil
}

func (c *Config) SearchPlayer(prefix string) (data.WGAccountList, error) {
	return c.wargaming.AccountListForSearch(prefix)
}

func (c *Config) OpenDirectory(path string) error {
	err := c.OpenWithDefaultApp(path)
	if err != nil {
		return failure.New(apperr.OpenDirectoryError, failure.Context{"path": path}, failure.Messagef("%s", err.Error()))
	}

	return nil
}

func (c *Config) GetUserConfig() (data.UserConfig, error) {
	cfg, err := c.persistence.LoadUserConfig()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return defaultUserConfig(), nil
		}

		return data.UserConfig{}, failure.Wrap(err)
	}

	return cfg, nil
}

func (c *Config) SaveUserConfig(cfg data.UserConfig) error {
	if err := c.persistence.SaveUserConfig(cfg); err != nil {
		return failure.Wrap(err)
	}

	return nil
}

func defaultUserConfig() data.UserConfig {
	return data.UserConfig{
		Version:      1,
		InstallPath:  "",
		ZoomRate:     100,
		StatsExtra:   "pvp_all",
		IsSendReport: true,
		Column: data.ColumnConfig{
			Player: data.PlayerColumnConfig{
				EnableNationFlag: false,
				ColorPattern:     "none",
			},
			Ship: data.ShipColumnConfig{
				EnableNationFlag: true,
				IsColored:        false,
			},
			Stats: data.StatsColumnConfig{
				Battles: data.DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         0,
				},
				Damage: data.DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         0,
				},
				MaxDamage: data.DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         0,
				},
				WinRate: data.DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         1,
				},
				SurvivedRate: data.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				KdRate: data.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         2,
				},
				Kill: data.DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         2,
				},
				Exp: data.DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         0,
				},
				PR: data.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         0,
				},
				HitRate: data.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				PlanesKilled: data.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         0,
				},
				PlatoonRate: data.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				EfficiencyBadge: data.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         0,
				},
				ThreatLevel: data.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				AvgTier: data.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				UsingShipTypeRate: data.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				UsingTierRate: data.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
			},
		},
	}
}
