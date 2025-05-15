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

	// file.
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
	//nolint:errcheck
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
		return readJSON(paths[0], data.TempArenaInfo{})
	}

	var latest data.TempArenaInfo
	for _, path := range paths {
		tempArenaInfo, err := readJSON(path, data.TempArenaInfo{})
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
