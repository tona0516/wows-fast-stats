package infra

import (
	"os"
	"path/filepath"
	"wfs/backend/apperr"
	"wfs/backend/vo"

	"github.com/pkg/errors"
)

const (
	ConfigDirName         string = "config"
	ConfigUserName        string = "user.json"
	ConfigAppName         string = "app.json"
	ConfigAlertPlayerName string = "alert_player.json"
)

//nolint:gochecknoglobals
var DefaultUserConfig vo.UserConfig = vo.UserConfig{
	FontSize:     "medium",
	SendReport:   true,
	StatsPattern: vo.StatsPatternPvPAll,
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
	CustomColor: vo.CustomColor{
		Skill: vo.SkillColor{
			Text: vo.SkillColorCode{
				Bad:         "#ff382d",
				BelowAvg:    "#fd9234",
				Avg:         "#ffd351",
				Good:        "#57e500",
				VeryGood:    "#44b200",
				Great:       "#02f7da",
				Unicum:      "#da6ff5",
				SuperUnicum: "#bf15ee",
			},
			Background: vo.SkillColorCode{
				Bad:         "#a41200",
				BelowAvg:    "#a34a02",
				Avg:         "#a38204",
				Good:        "#518517",
				VeryGood:    "#2f6f41",
				Great:       "#04436d",
				Unicum:      "#232166",
				SuperUnicum: "#531460",
			},
		},
		Tier: vo.TierColor{
			Own: vo.TierColorCode{
				Low:    "#8CA113",
				Middle: "#205B85",
				High:   "#990F4F",
			},
			Other: vo.TierColorCode{
				Low:    "#E6F5B0",
				Middle: "#B3D7DD",
				High:   "#E3ADD5",
			},
		},
		ShipType: vo.ShipTypeColor{
			Own: vo.ShipTypeColorCode{
				CV: "#5E2883",
				BB: "#CA1028",
				CL: "#27853F",
				DD: "#D9760F",
				SS: "#233B8B",
			},
			Other: vo.ShipTypeColorCode{
				CV: "#CAB2D6",
				BB: "#FBB4C4",
				CL: "#CCEBC5",
				DD: "#FEE6AA",
				SS: "#B3CDE3",
			},
		},
	},
	CustomDigit: vo.CustomDigit{
		PR:                0,
		Damage:            0,
		WinRate:           1,
		KdRate:            2,
		Kill:              2,
		PlanesKilled:      1,
		Exp:               0,
		Battles:           0,
		SurvivedRate:      1,
		HitRate:           1,
		AvgTier:           2,
		UsingShipTypeRate: 1,
		UsingTierRate:     1,
	},
}

type Config struct{}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) User() (vo.UserConfig, error) {
	// note: set default value
	return read(ConfigUserName, DefaultUserConfig)
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
	result, err := readJSON(path, defaultValue)
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
