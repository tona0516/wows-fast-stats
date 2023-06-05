package infra

import (
	"changeme/backend/apperr"
	"changeme/backend/vo"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	ConfigDirName         string = "config"
	ConfigUserName        string = "user.json"
	ConfigAppName         string = "app.json"
	ConfigAlertPlayerName string = "alert_player.json"
)

//nolint:gochecknoglobals
var defaultUserConfig vo.UserConfig = vo.UserConfig{
	FontSize: "medium",
	Displays: vo.Displays{
		Basic: vo.Basic{
			IsInAvg:    true,
			PlayerName: true,
			ShipInfo:   true,
		},
		Ship: vo.Ship{
			PR:           true,
			Damage:       true,
			WinRate:      true,
			KdRate:       false,
			Exp:          false,
			Battles:      true,
			SurvivedRate: false,
			HitRate:      false,
		},
		Overall: vo.Overall{
			Damage:            true,
			WinRate:           true,
			KdRate:            false,
			Exp:               false,
			Battles:           true,
			SurvivedRate:      false,
			AvgTier:           false,
			UsingShipTypeRate: false,
			UsingTierRate:     false,
		},
	},
	StatsPattern: vo.StatsPatternPvPAll,
}

type Config struct{}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) User() (vo.UserConfig, error) {
	// note: set default value
	return read(ConfigUserName, defaultUserConfig)
}

func (c *Config) UpdateUser(config vo.UserConfig) error {
	return update(ConfigUserName, config)
}

func (c *Config) App() (vo.AppConfig, error) {
	return read(ConfigAppName, vo.AppConfig{})
}

func (c *Config) UpdateApp(config vo.AppConfig) error {
	return update(ConfigAppName, config)
}

func (c *Config) AlertPlayers() ([]vo.AlertPlayer, error) {
	return read(ConfigAlertPlayerName, make([]vo.AlertPlayer, 0))
}

func (c *Config) UpdateAlertPlayer(player vo.AlertPlayer) error {
	players, err := read(ConfigAlertPlayerName, make([]vo.AlertPlayer, 0))
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

	return update(ConfigAlertPlayerName, players)
}

func (c *Config) RemoveAlertPlayer(accountID int) error {
	players, err := read(ConfigAlertPlayerName, make([]vo.AlertPlayer, 0))
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

	return update(ConfigAlertPlayerName, players)
}

func read[T any](filename string, defaultValue T) (T, error) {
	errDetail := apperr.Cfg.Read

	_ = os.Mkdir(ConfigDirName, 0o755)

	path := filepath.Join(ConfigDirName, filename)
	if _, err := os.Stat(path); err != nil {
		//nolint:nilerr
		return defaultValue, nil
	}

	f, err := os.ReadFile(path)
	if err != nil {
		return defaultValue, errors.WithStack(errDetail.WithRaw(err))
	}
	if err := json.Unmarshal(f, &defaultValue); err != nil {
		return defaultValue, errors.WithStack(errDetail.WithRaw(err))
	}

	return defaultValue, nil
}

func update[T any](filename string, target T) error {
	errDetail := apperr.Cfg.Update

	_ = os.Mkdir(ConfigDirName, 0o755)

	file, err := os.Create(filepath.Join(ConfigDirName, filename))
	if err != nil {
		return errors.WithStack(errDetail.WithRaw(err))
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(target); err != nil {
		return errors.WithStack(errDetail.WithRaw(err))
	}

	return nil
}
