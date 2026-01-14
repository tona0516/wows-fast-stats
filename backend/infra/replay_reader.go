package infra

import (
	"io/fs"
	"os"
	"path/filepath"
	"wfs/backend/data"

	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
)

type ReplayReader struct {
	replayDir string
	fileName  string
}

func NewReplayReader(i do.Injector) (*ReplayReader, error) {
	return &ReplayReader{
		replayDir: "replays",
		fileName:  "tempArenaInfo.json",
	}, nil
}

func (r *ReplayReader) TempArenaInfo(installPath string) (data.TempArenaInfo, error) {
	var tempArenaInfo data.TempArenaInfo

	tempArenaInfoPaths := []string{}
	root := filepath.Join(installPath, r.replayDir)
	if _, err := os.Stat(root); err != nil {
		return tempArenaInfo, failure.Translate(err, data.ErrTempArenaInfoNotFound)
	}

	err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.Name() != r.fileName {
			return nil
		}

		tempArenaInfoPaths = append(tempArenaInfoPaths, path)
		return nil
	})
	if err != nil {
		return tempArenaInfo, failure.Translate(err, data.ErrTempArenaInfoSearch)
	}

	return r.decideTempArenaInfo(tempArenaInfoPaths)
}

func (r *ReplayReader) decideTempArenaInfo(paths []string) (data.TempArenaInfo, error) {
	var result data.TempArenaInfo
	size := len(paths)

	if size == 0 {
		return result, failure.New(data.ErrTempArenaInfoNotFound)
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
		return result, failure.New(data.ErrTempArenaInfoNotFound)
	}

	return latest, nil
}
