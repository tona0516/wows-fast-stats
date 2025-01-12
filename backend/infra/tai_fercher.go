package infra

import (
	"io/fs"
	"os"
	"path/filepath"
	"wfs/backend/apperr"
	"wfs/backend/domain/model"

	"github.com/morikuni/failure"
)

const (
	replaysDir        string = "replays"
	tempArenaInfoFile string = "tempArenaInfo.json"
)

type TAIFetcher struct{}

func NewTaiFetcher() *TAIFetcher {
	return &TAIFetcher{}
}

func (f *TAIFetcher) Get(installPath string) (model.TempArenaInfo, error) {
	var result model.TempArenaInfo

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

	return f.decide(paths)
}

func (f *TAIFetcher) decide(paths []string) (model.TempArenaInfo, error) {
	var result model.TempArenaInfo
	size := len(paths)

	if size == 0 {
		return result, failure.New(apperr.FileNotExist)
	}

	if size == 1 {
		return readJSON(paths[0], model.TempArenaInfo{})
	}

	var latest model.TempArenaInfo
	for _, path := range paths {
		tempArenaInfo, err := readJSON(path, model.TempArenaInfo{})
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
