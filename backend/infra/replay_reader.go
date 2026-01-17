package infra

import (
	"io/fs"
	"os"
	"path/filepath"
	"wfs/backend/core"

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

func (r *ReplayReader) TempArenaInfo(installPath string) (core.TempArenaInfo, error) {
	var tempArenaInfo core.TempArenaInfo

	tempArenaInfoPaths := []string{}
	root := filepath.Join(installPath, r.replayDir)
	if _, err := os.Stat(root); err != nil {
		return tempArenaInfo, failure.Translate(err, core.ErrTempArenaInfoNotFound)
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
		return tempArenaInfo, failure.Translate(err, core.ErrTempArenaInfoSearch)
	}

	return r.decideTempArenaInfo(tempArenaInfoPaths)
}

func (r *ReplayReader) decideTempArenaInfo(paths []string) (core.TempArenaInfo, error) {
	var result core.TempArenaInfo
	size := len(paths)

	if size == 0 {
		return result, failure.New(core.ErrTempArenaInfoNotFound)
	}

	if size == 1 {
		return readJSON[core.TempArenaInfo](paths[0])
	}

	var latest core.TempArenaInfo
	for _, path := range paths {
		tempArenaInfo, err := readJSON[core.TempArenaInfo](path)
		if err != nil {
			continue
		}

		if tempArenaInfo.Unixtime() > latest.Unixtime() {
			latest = tempArenaInfo
		}
	}

	if latest.Unixtime() == 0 {
		return result, failure.New(core.ErrTempArenaInfoNotFound)
	}

	return latest, nil
}
