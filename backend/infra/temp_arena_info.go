package infra

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"wfs/backend/apperr"
	"wfs/backend/vo"

	"github.com/pkg/errors"
)

const (
	tempArenaInfoDir  string = "temp_arena_info"
	ReplayDir         string = "replays"
	TempArenaInfoName string = "tempArenaInfo.json"
)

var errNoTempArenaInfo = errors.New("no tempArenaInfo.json")

type TempArenaInfo struct{}

func NewTempArenaInfo() *TempArenaInfo {
	return &TempArenaInfo{}
}

func (t *TempArenaInfo) Get(installPath string) (vo.TempArenaInfo, error) {
	var tempArenaInfo vo.TempArenaInfo

	tempArenaInfoPaths := []string{}
	root := filepath.Join(installPath, ReplayDir)
	err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.Name() != TempArenaInfoName {
			return nil
		}

		tempArenaInfoPaths = append(tempArenaInfoPaths, path)
		return nil
	})

	if err != nil {
		return tempArenaInfo, apperr.New(apperr.ReadFile, err)
	}

	tempArenaInfo, err = decideTempArenaInfo(tempArenaInfoPaths)
	if err != nil {
		return tempArenaInfo, err
	}

	return tempArenaInfo, nil
}

func (t *TempArenaInfo) Save(tempArenaInfo vo.TempArenaInfo) error {
	_ = os.Mkdir(tempArenaInfoDir, 0o755)
	path := filepath.Join(tempArenaInfoDir, "tempArenaInfo_"+strconv.FormatInt(tempArenaInfo.Unixtime(), 10)+".json")

	if err := writeJSON(path, tempArenaInfo); err != nil {
		return apperr.New(apperr.WriteFile, err)
	}

	return nil
}

func decideTempArenaInfo(paths []string) (vo.TempArenaInfo, error) {
	size := len(paths)

	if size == 0 {
		return vo.TempArenaInfo{}, apperr.New(apperr.ReadFile, errNoTempArenaInfo)
	}

	if size == 1 {
		tempArenaInfo, err := readJSON(paths[0], vo.TempArenaInfo{})
		if err != nil {
			return vo.TempArenaInfo{}, apperr.New(apperr.ReadFile, err)
		}

		return tempArenaInfo, nil
	}

	var latest *vo.TempArenaInfo
	var latestUnixtime int64
	for _, path := range paths {
		tempArenaInfo, err := readJSON(path, vo.TempArenaInfo{})
		if err != nil {
			continue
		}

		unixtime := tempArenaInfo.Unixtime()

		if unixtime > latestUnixtime {
			latest = &tempArenaInfo
			latestUnixtime = unixtime
		}
	}

	if latest == nil {
		return vo.TempArenaInfo{}, apperr.New(apperr.ReadFile, errNoTempArenaInfo)
	}

	return *latest, nil
}
