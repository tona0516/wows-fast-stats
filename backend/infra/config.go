package infra

import (
	"changeme/backend/apperr"
	"changeme/backend/vo"
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
			PR:      true,
			Damage:  true,
			WinRate: true,
			Battles: true,
		},
		Overall: vo.Overall{
			Damage:  true,
			WinRate: true,
			Battles: true,
		},
	},
	StatsPattern: vo.StatsPatternPvPAll,
	SendReport:   true,
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
	path := filepath.Join(ConfigDirName, filename)
	result, err := readJSON[T](path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return defaultValue, nil
		}

		return defaultValue, apperr.New(apperr.ReadFile, err)
	}

	return result, nil
}

func update[T any](filename string, target T) error {
	path := filepath.Join(ConfigDirName, filename)
	err := writeJSON(path, target)
	if err != nil {
		return apperr.New(apperr.WriteFile, err)
	}

	return nil
}
