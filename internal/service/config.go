package service

import (
	"errors"
	"os"
	"path/filepath"
	"wfs/internal/apperr"
	"wfs/internal/data"
	"wfs/internal/infra"

	"github.com/morikuni/failure"
)

const GameExeName = "WorldOfWarships.exe"

type Config struct {
	localStorage infra.LocalStorage
}

func NewConfig(
	localStorage infra.LocalStorage,
) *Config {
	return &Config{
		localStorage: localStorage,
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

func (c *Config) GetUserConfig() (data.UserConfig, error) {
	cfg, err := c.localStorage.UserConfig()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return defaultUserConfig(), nil
		}

		return data.UserConfig{}, failure.Wrap(err)
	}

	return cfg, nil
}

func (c *Config) SaveUserConfig(cfg data.UserConfig) error {
	if err := c.localStorage.SetUserConfig(cfg); err != nil {
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
