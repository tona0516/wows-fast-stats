package persistence

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/morikuni/failure"
)

const (
	appDir               = "wows-fast-stats"
	userSettingFileName  = "user_setting.json"
	alertPlayersFileName = "alert_players.json"
)

type API interface {
	GetUserSetting() (UserSetting, error)
	SaveUserSetting(us UserSetting) error
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

func (a *api) SaveUserSetting(us UserSetting) error {
	path, err := a.getPath(userSettingFileName)
	if err != nil {
		return failure.Wrap(err)
	}

	return writeJSON(path, us)
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
