package service

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/domain"
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

func (c *Config) GetUserConfig() (domain.UserConfig, error) {
	cfg, err := c.persistence.LoadUserConfig()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return defaultUserConfig(), nil
		}

		return domain.UserConfig{}, failure.Wrap(err)
	}

	return cfg, nil
}

func (c *Config) SaveUserConfig(cfg domain.UserConfig) error {
	if err := c.persistence.SaveUserConfig(cfg); err != nil {
		return failure.Wrap(err)
	}

	return nil
}

func defaultUserConfig() domain.UserConfig {
	return domain.UserConfig{
		Version:      1,
		InstallPath:  "",
		ZoomRate:     100,
		StatsExtra:   "pvp_all",
		IsSendReport: true,
		Column: domain.ColumnConfig{
			Player: domain.PlayerColumnConfig{
				EnableNationFlag: false,
				ColorPattern:     "none",
			},
			Ship: domain.ShipColumnConfig{
				EnableNationFlag: true,
				IsColored:        false,
			},
			Stats: domain.StatsColumnConfig{
				Battles: domain.DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         0,
				},
				Damage: domain.DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         0,
				},
				MaxDamage: domain.DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         0,
				},
				WinRate: domain.DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         1,
				},
				SurvivedRate: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				KdRate: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         2,
				},
				Kill: domain.DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         2,
				},
				Exp: domain.DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         0,
				},
				PR: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         0,
				},
				HitRate: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				PlanesKilled: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         0,
				},
				PlatoonRate: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				EfficiencyBadge: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         0,
				},
				ThreatLevel: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				AvgTier: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				UsingShipTypeRate: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				UsingTierRate: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
			},
		},
	}
}
