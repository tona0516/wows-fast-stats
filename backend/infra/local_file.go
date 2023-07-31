package infra

import (
	"encoding/base64"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"wfs/backend/apperr"
	"wfs/backend/application/vo"
	"wfs/backend/domain"

	"github.com/pkg/errors"
)

const (
	// directory.
	configDir        string = "config"
	replaysDir       string = "replays"
	tempArenaInfoDir string = "temp_arena_info"
	// file.
	userConfigFile    string = "user.json"
	appConfigFile     string = "app.json"
	alertPlayerFile   string = "alert_player.json"
	tempArenaInfoFile string = "tempArenaInfo.json"
)

//nolint:gochecknoglobals
var DefaultUserConfig domain.UserConfig = domain.UserConfig{
	FontSize:        "medium",
	SendReport:      true,
	NotifyUpdatable: true,
	StatsPattern:    domain.StatsPatternPvPAll,
	Displays: domain.Displays{
		Basic: domain.Basic{
			IsInAvg:    true,
			PlayerName: true,
			ShipInfo:   true,
		},
		Ship: domain.Ship{
			PR:      true,
			Damage:  true,
			WinRate: true,
			Battles: true,
		},
		Overall: domain.Overall{
			Damage:  true,
			WinRate: true,
			Battles: true,
		},
	},
	CustomColor: domain.CustomColor{
		Skill: domain.SkillColor{
			Text: domain.SkillColorCode{
				Bad:         "#ff382d",
				BelowAvg:    "#fd9234",
				Avg:         "#ffd351",
				Good:        "#57e500",
				VeryGood:    "#44b200",
				Great:       "#02f7da",
				Unicum:      "#da6ff5",
				SuperUnicum: "#bf15ee",
			},
			Background: domain.SkillColorCode{
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
		Tier: domain.TierColor{
			Own: domain.TierColorCode{
				Low:    "#8CA113",
				Middle: "#205B85",
				High:   "#990F4F",
			},
			Other: domain.TierColorCode{
				Low:    "#E6F5B0",
				Middle: "#B3D7DD",
				High:   "#E3ADD5",
			},
		},
		ShipType: domain.ShipTypeColor{
			Own: domain.ShipTypeColorCode{
				CV: "#5E2883",
				BB: "#CA1028",
				CL: "#27853F",
				DD: "#D9760F",
				SS: "#233B8B",
			},
			Other: domain.ShipTypeColorCode{
				CV: "#CAB2D6",
				BB: "#FBB4C4",
				CL: "#CCEBC5",
				DD: "#FEE6AA",
				SS: "#B3CDE3",
			},
		},
	},
	CustomDigit: domain.CustomDigit{
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
	TeamAverage: domain.TeamAverage{
		MinShipBattles:    1,
		MinOverallBattles: 10,
	},
}

type LocalFile struct {
	userConfigPath  string
	appConfigPath   string
	alertPlayerPath string
}

func NewLocalFile() *LocalFile {
	return &LocalFile{
		userConfigPath:  filepath.Join(configDir, userConfigFile),
		appConfigPath:   filepath.Join(configDir, appConfigFile),
		alertPlayerPath: filepath.Join(configDir, alertPlayerFile),
	}
}

func (l *LocalFile) User() (domain.UserConfig, error) {
	config, err := readJSON[domain.UserConfig](l.userConfigPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return DefaultUserConfig, nil
		}
		return domain.UserConfig{}, apperr.New(apperr.ReadFile, err)
	}

	return config, nil
}

func (l *LocalFile) UpdateUser(config domain.UserConfig) error {
	if err := writeJSON(l.userConfigPath, config); err != nil {
		return apperr.New(apperr.WriteFile, err)
	}

	return nil
}

func (l *LocalFile) App() (vo.AppConfig, error) {
	config, err := readJSON[vo.AppConfig](l.appConfigPath)
	if err != nil {
		config = vo.AppConfig{}
		if errors.Is(err, fs.ErrNotExist) {
			return config, nil
		}
		return config, apperr.New(apperr.ReadFile, err)
	}

	return config, nil
}

func (l *LocalFile) UpdateApp(config vo.AppConfig) error {
	if err := writeJSON(l.appConfigPath, config); err != nil {
		return apperr.New(apperr.WriteFile, err)
	}

	return nil
}

func (l *LocalFile) AlertPlayers() ([]domain.AlertPlayer, error) {
	players, err := readJSON[[]domain.AlertPlayer](l.alertPlayerPath)
	if err != nil {
		players = make([]domain.AlertPlayer, 0)
		if errors.Is(err, fs.ErrNotExist) {
			return players, nil
		}
		return players, apperr.New(apperr.ReadFile, err)
	}

	return players, nil
}

func (l *LocalFile) UpdateAlertPlayer(player domain.AlertPlayer) error {
	players, err := l.AlertPlayers()
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

	if err := writeJSON(l.alertPlayerPath, players); err != nil {
		return apperr.New(apperr.WriteFile, err)
	}

	return nil
}

func (l *LocalFile) RemoveAlertPlayer(accountID int) error {
	players, err := l.AlertPlayers()
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

	if err := writeJSON(l.alertPlayerPath, players); err != nil {
		return apperr.New(apperr.WriteFile, err)
	}

	return nil
}

func (l *LocalFile) SaveScreenshot(path string, base64Data string) error {
	dir := filepath.Dir(path)
	_ = os.Mkdir(dir, 0o755)

	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return apperr.New(apperr.DecodeBase64, err)
	}

	f, err := os.Create(path)
	if err != nil {
		return apperr.New(apperr.WriteFile, err)
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return apperr.New(apperr.WriteFile, err)
	}

	return nil
}

func (l *LocalFile) TempArenaInfo(installPath string) (domain.TempArenaInfo, error) {
	var tempArenaInfo domain.TempArenaInfo

	tempArenaInfoPaths := []string{}
	root := filepath.Join(installPath, replaysDir)
	err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.Name() != tempArenaInfoFile {
			return nil
		}

		tempArenaInfoPaths = append(tempArenaInfoPaths, path)
		return nil
	})

	if err != nil {
		return tempArenaInfo, apperr.New(apperr.ReadFile, err)
	}

	tempArenaInfo, err = decideTempArenaInfo(tempArenaInfoPaths)
	if err != nil {
		return tempArenaInfo, apperr.New(apperr.ReadFile, err)
	}

	return tempArenaInfo, nil
}

func (l *LocalFile) SaveTempArenaInfo(tempArenaInfo domain.TempArenaInfo) error {
	path := filepath.Join(tempArenaInfoDir, "tempArenaInfo_"+strconv.FormatInt(tempArenaInfo.Unixtime(), 10)+".json")

	if err := writeJSON(path, tempArenaInfo); err != nil {
		return apperr.New(apperr.WriteFile, err)
	}

	return nil
}

func decideTempArenaInfo(paths []string) (domain.TempArenaInfo, error) {
	var result domain.TempArenaInfo
	size := len(paths)

	if size == 0 {
		return result, fs.ErrNotExist
	}

	if size == 1 {
		result, err := readJSON[domain.TempArenaInfo](paths[0])
		if err != nil {
			return result, err
		}

		return result, nil
	}

	var latest domain.TempArenaInfo
	for _, path := range paths {
		tempArenaInfo, err := readJSON[domain.TempArenaInfo](path)
		if err != nil {
			continue
		}

		if tempArenaInfo.Unixtime() > latest.Unixtime() {
			latest = tempArenaInfo
		}
	}

	if latest.Unixtime() == 0 {
		return result, fs.ErrNotExist
	}

	return latest, nil
}

func readJSON[T any](path string) (T, error) {
	var result T

	f, err := os.ReadFile(path)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(f, &result); err != nil {
		return result, err
	}

	return result, nil
}

func writeJSON[T any](path string, target T) error {
	_ = os.MkdirAll(filepath.Dir(path), 0o755)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")

	return encoder.Encode(target)
}
