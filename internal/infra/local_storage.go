package infra

import (
	"io/fs"
	"os"
	"path/filepath"
	"wfs/internal/data"

	"github.com/morikuni/failure"
)

const (
	replaysDir        string = "replays"
	tempArenaInfoFile string = "tempArenaInfo.json"
	ownIGNFile        string = "own_ign.txt"
	expectedStatsFile string = "expected_stats.json"
	userConfigFile    string = "user_config.json"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type LocalStorage interface {
	TempArenaInfo(installPath string) (data.TempArenaInfo, error)
	OwnIGN() (string, error)
	ExpectedStats() (data.NSExpectedStats, error)
	UserConfig() (data.UserConfig, error)
	SetOwnIGN(ign string) error
	SetExpectedStats(data data.NSExpectedStats) error
	SetUserConfig(data data.UserConfig) error
}

type localStorage struct {
	userDataDir string
}

func NewLocalStorage(userDataDir string) LocalStorage {
	return &localStorage{userDataDir: userDataDir}
}

func (s *localStorage) TempArenaInfo(installPath string) (data.TempArenaInfo, error) {
	var tempArenaInfo data.TempArenaInfo

	tempArenaInfoPaths := []string{}
	root := filepath.Join(installPath, replaysDir)
	if _, err := os.Stat(root); err != nil {
		return tempArenaInfo, failure.Wrap(err)
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

	return s.decideTempArenaInfo(tempArenaInfoPaths)
}

func (s *localStorage) OwnIGN() (string, error) {
	return readString(filepath.Join(s.userDataDir, ownIGNFile))
}

func (s *localStorage) ExpectedStats() (data.NSExpectedStats, error) {
	return readJSON[data.NSExpectedStats](filepath.Join(s.userDataDir, expectedStatsFile))
}

func (s *localStorage) UserConfig() (data.UserConfig, error) {
	return readJSON[data.UserConfig](filepath.Join(s.userDataDir, userConfigFile))
}

func (s *localStorage) SetOwnIGN(ign string) error {
	return writeString(filepath.Join(s.userDataDir, ownIGNFile), ign)
}

func (s *localStorage) SetExpectedStats(data data.NSExpectedStats) error {
	return writeJSON(filepath.Join(s.userDataDir, expectedStatsFile), data)
}

func (s *localStorage) SetUserConfig(data data.UserConfig) error {
	return writeJSON(filepath.Join(s.userDataDir, userConfigFile), data)
}

func (s *localStorage) decideTempArenaInfo(paths []string) (data.TempArenaInfo, error) {
	var result data.TempArenaInfo
	size := len(paths)

	if size == 0 {
		return result, failure.Wrap(fs.ErrNotExist)
	}

	if size == 1 {
		return readJSON[data.TempArenaInfo](paths[0])
	}

	var latest data.TempArenaInfo
	for _, path := range paths {
		tempArenaInfo, err := readJSON[data.TempArenaInfo](path)
		if err != nil {
			continue
		}

		if tempArenaInfo.Unixtime() > latest.Unixtime() {
			latest = tempArenaInfo
		}
	}

	if latest.Unixtime() == 0 {
		return result, failure.Wrap(fs.ErrNotExist)
	}

	return latest, nil
}
