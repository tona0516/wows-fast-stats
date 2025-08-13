package infra

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
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

type LocalFile struct{}

func NewLocalFile() *LocalFile {
	return &LocalFile{}
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

	return l.decideTempArenaInfo(tempArenaInfoPaths)
}

func (l *LocalFile) decideTempArenaInfo(paths []string) (data.TempArenaInfo, error) {
	var result data.TempArenaInfo
	size := len(paths)

	if size == 0 {
		return result, failure.New(apperr.FileNotExist)
	}

	if size == 1 {
		return l.read(paths[0])
	}

	var latest data.TempArenaInfo
	for _, path := range paths {
		tempArenaInfo, err := l.read(path)
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

func (l *LocalFile) read(path string) (data.TempArenaInfo, error) {
	var result data.TempArenaInfo

	f, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return result, failure.New(apperr.FileNotExist)
		}
		return result, failure.Wrap(err)
	}

	if err = json.Unmarshal(f, &result); err != nil {
		return result, failure.Wrap(err)
	}

	return result, nil
}
