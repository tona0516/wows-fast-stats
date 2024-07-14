package infra

import (
	"encoding/base64"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"wfs/backend/apperr"
	"wfs/backend/data"

	"github.com/morikuni/failure"
)

const (
	// directory.
	replaysDir       string = "replays"
	tempArenaInfoDir string = "temp_arena_info"
	configDir        string = "config_v3"
	cacheDir         string = "cache_v2"

	// file.
	tempArenaInfoFile string = "tempArenaInfo.json"
	ignFile           string = "ign.txt"
	expectedStatsFile string = "expected_stats.json"
	userConfigFile    string = "user_config_v2.json"
	alertPlayerFile   string = "alert_player_v1.json"
)

type LocalFile struct {
	ignPath           string
	expectedStatsPath string
	userConfigPath    string
	alertPlayerPath   string
}

func NewLocalFile() *LocalFile {
	return &LocalFile{
		ignPath:           filepath.Join(cacheDir, ignFile),
		expectedStatsPath: filepath.Join(cacheDir, expectedStatsFile),
		userConfigPath:    filepath.Join(configDir, userConfigFile),
		alertPlayerPath:   filepath.Join(configDir, alertPlayerFile),
	}
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

func (l *LocalFile) TempArenaInfo(installPath string) (data.TempArenaInfo, error) {
	var tempArenaInfo data.TempArenaInfo

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

func (l *LocalFile) SaveTempArenaInfo(tempArenaInfo data.TempArenaInfo) error {
	path := filepath.Join(tempArenaInfoDir, "tempArenaInfo_"+strconv.FormatInt(tempArenaInfo.Unixtime(), 10)+".json")
	return writeJSON(path, tempArenaInfo)
}

func decideTempArenaInfo(paths []string) (data.TempArenaInfo, error) {
	var result data.TempArenaInfo
	size := len(paths)

	if size == 0 {
		return result, failure.New(apperr.FileNotExist)
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
		return result, failure.New(apperr.FileNotExist)
	}

	return latest, nil
}

func (l *LocalFile) IGN() (string, error) {
	b, err := readFile(l.ignPath)
	return string(b), err
}

func (l *LocalFile) WriteIGN(ign string) error {
	return writeFile(l.ignPath, []byte(ign))
}

func (l *LocalFile) ExpectedStats() (data.ExpectedStats, error) {
	return readJSON[data.ExpectedStats](l.expectedStatsPath)
}

func (l *LocalFile) WriteExpectedStats(target data.ExpectedStats) error {
	return writeJSON(l.expectedStatsPath, target)
}

func (l *LocalFile) UserConfigV2() (data.UserConfigV2, error) {
	return readJSON[data.UserConfigV2](l.userConfigPath)
}

func (l *LocalFile) WriteUserConfigV2(target data.UserConfigV2) error {
	return writeJSON(l.userConfigPath, target)
}

func (l *LocalFile) AlertPlayerV1() ([]data.AlertPlayer, error) {
	return readJSON[[]data.AlertPlayer](l.alertPlayerPath)
}

func (l *LocalFile) WriteAlertPlayerV1(target []data.AlertPlayer) error {
	return writeJSON(l.alertPlayerPath, target)
}
