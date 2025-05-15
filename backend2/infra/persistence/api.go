package persistence

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"wfs/backend/apperr"

	"github.com/morikuni/failure"
)

const (
	appDir                      = "wows-fast-stats"
	userSettingFileName         = "user_setting.json"
	alertPlayersFileName        = "alert_players.json"
	tempArenaInfoFile    string = "tempArenaInfo.json"
	replaysDir           string = "replays"
)

type API interface {
	GetUserSetting() (UserSetting, error)
	GetTempArenaInfo(installPath string) (TempArenaInfo, error)
	SaveUserSetting(us UserSetting) error
	SaveTempArenaInfo(path string, tempArenaInfo TempArenaInfo) error
	SaveScreenshot(path string, base64Data string) error
}

type GetConfigDirFunc func() (string, error)

type api struct {
	getConfigDir GetConfigDirFunc
}

func NewAPI(getConfigDir GetConfigDirFunc) API {
	return &api{getConfigDir: getConfigDir}
}

func (a *api) GetUserSetting() (UserSetting, error) {
	path, err := a.getPath(userSettingFileName)
	if err != nil {
		return UserSetting{}, failure.Wrap(err)
	}

	return readJSON[UserSetting](path)
}

func (a *api) GetTempArenaInfo(installPath string) (TempArenaInfo, error) {
	var result TempArenaInfo

	root := filepath.Join(installPath, replaysDir)
	if _, err := os.Stat(root); err != nil {
		return result, failure.Translate(err, apperr.ReplayDirNotFoundError)
	}

	paths := make([]string, 0)

	err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return failure.Wrap(err)
		}

		if info.IsDir() {
			return nil
		}

		if info.Name() != tempArenaInfoFile {
			return nil
		}

		paths = append(paths, path)

		return nil
	})
	if err != nil {
		return result, err
	}

	size := len(paths)
	if size == 0 {
		return result, failure.New(apperr.FileNotExist)
	}

	if size == 1 {
		return readJSON[TempArenaInfo](paths[0])
	}

	var latest TempArenaInfo

	for _, path := range paths {
		tempArenaInfo, err := readJSON[TempArenaInfo](path)
		if err != nil {
			continue
		}

		if tempArenaInfo.Unixtime() > latest.Unixtime() {
			latest = tempArenaInfo
		}
	}

	if latest.Unixtime() == 0 {
		return result, failure.New(apperr.FileNotExist)
	}

	return latest, nil
}

func (a *api) SaveUserSetting(us UserSetting) error {
	path, err := a.getPath(userSettingFileName)
	if err != nil {
		return failure.Wrap(err)
	}

	return writeJSON(path, us)
}

func (a *api) SaveTempArenaInfo(path string, tempArenaInfo TempArenaInfo) error {
	// TODO
	return nil
}

func (a *api) SaveScreenshot(path string, base64Data string) error {
	// TODO
	return nil
}

func (a *api) getPath(name string) (string, error) {
	dir, err := a.getConfigDir()
	if err != nil {
		return "", failure.Wrap(err)
	}

	return filepath.Join(dir, appDir, name), nil
}

func readJSON[T any](path string) (T, error) {
	var result T

	//nolint:gosec
	f, err := os.ReadFile(path)
	if err != nil {
		return result, failure.Wrap(err)
	}

	if err = json.Unmarshal(f, &result); err != nil {
		return result, failure.Wrap(err)
	}

	return result, nil
}

func writeJSON[T any](path string, data T) error {
	//nolint:gosec
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return failure.Wrap(err)
	}

	//nolint:gosec
	f, err := os.Create(path)
	if err != nil {
		return failure.Wrap(err)
	}
	//nolint:errcheck
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")

	if err = encoder.Encode(data); err != nil {
		return failure.Wrap(err)
	}

	return nil
}
