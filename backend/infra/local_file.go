package infra

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"wfs/backend/apperr"
	"wfs/backend/domain"

	"github.com/morikuni/failure"
)

const (
	// directory.
	ConfigDir        string = "config"
	replaysDir       string = "replays"
	tempArenaInfoDir string = "temp_arena_info"

	// file.
	UserConfigFile    string = "user.json"
	AlertPlayerFile   string = "alert_player.json"
	tempArenaInfoFile string = "tempArenaInfo.json"
)

type LocalFile struct {
	userConfigPath  string
	alertPlayerPath string
}

func NewLocalFile() *LocalFile {
	return &LocalFile{
		userConfigPath:  filepath.Join(ConfigDir, UserConfigFile),
		alertPlayerPath: filepath.Join(ConfigDir, AlertPlayerFile),
	}
}

func (l *LocalFile) User() (domain.UserConfig, error) {
	config, err := readJSON(l.userConfigPath, domain.DefaultUserConfig)
	if err != nil && failure.Is(err, apperr.FileNotExist) {
		return domain.DefaultUserConfig, nil
	}

	return config, err
}

func (l *LocalFile) AlertPlayers() ([]domain.AlertPlayer, error) {
	players, err := readJSON(l.alertPlayerPath, []domain.AlertPlayer{})
	if err != nil && failure.Is(err, apperr.FileNotExist) {
		return []domain.AlertPlayer{}, nil
	}

	return players, err
}

func (l *LocalFile) SaveScreenshot(path string, base64Data string) error {
	dir := filepath.Dir(path)
	_ = os.Mkdir(dir, 0o755)

	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return failure.Wrap(err)
	}

	f, err := os.Create(path)
	if err != nil {
		return failure.Wrap(err)
	}
	defer f.Close()

	_, err = f.Write(data)
	return failure.Wrap(err)
}

func (l *LocalFile) TempArenaInfo(installPath string) (domain.TempArenaInfo, error) {
	var tempArenaInfo domain.TempArenaInfo

	tempArenaInfoPaths := []string{}
	root := filepath.Join(installPath, replaysDir)
	if _, err := os.Stat(root); err != nil {
		return tempArenaInfo, failure.New(apperr.ReplayDirNotFoundError, failure.Messagef("%s", err.Error()))
	}

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
		return tempArenaInfo, failure.Wrap(err)
	}

	return decideTempArenaInfo(tempArenaInfoPaths)
}

func (l *LocalFile) SaveTempArenaInfo(tempArenaInfo domain.TempArenaInfo) error {
	path := filepath.Join(tempArenaInfoDir, "tempArenaInfo_"+strconv.FormatInt(tempArenaInfo.Unixtime(), 10)+".json")
	return writeJSON(path, tempArenaInfo)
}

func (l *LocalFile) IsExistUser() bool {
	_, err := os.Stat(l.userConfigPath)
	return err == nil
}

func (l *LocalFile) DeleteUser() error {
	return os.RemoveAll(l.userConfigPath)
}

func (l *LocalFile) IsExistAlertPlayers() bool {
	_, err := os.Stat(l.alertPlayerPath)
	return err == nil
}

func (l *LocalFile) DeleteAlertPlayers() error {
	return os.RemoveAll(l.alertPlayerPath)
}

func decideTempArenaInfo(paths []string) (domain.TempArenaInfo, error) {
	var result domain.TempArenaInfo
	size := len(paths)

	if size == 0 {
		return result, failure.New(apperr.FileNotExist)
	}

	if size == 1 {
		return readJSON(paths[0], domain.TempArenaInfo{})
	}

	var latest domain.TempArenaInfo
	for _, path := range paths {
		tempArenaInfo, err := readJSON(path, domain.TempArenaInfo{})
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

func readJSON[T any](path string, defaulValue T) (T, error) {
	errCtx := failure.Context{"path": path}

	f, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return defaulValue, failure.New(apperr.FileNotExist, errCtx)
		}
		return defaulValue, failure.Wrap(err, errCtx)
	}
	errCtx["target"] = string(f)

	err = json.Unmarshal(f, &defaulValue)
	return defaulValue, failure.Wrap(err, errCtx)
}

func writeJSON[T any](path string, target T) error {
	//nolint:errchkjson
	marshaled, _ := json.Marshal(target)
	errCtx := failure.Context{"path": path, "target": string(marshaled)}

	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	f, err := os.Create(path)
	if err != nil {
		return failure.Wrap(err, errCtx)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(target)
	return failure.Wrap(err, errCtx)
}
