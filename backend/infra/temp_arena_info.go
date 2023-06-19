package infra

import (
	"changeme/backend/apperr"
	"changeme/backend/vo"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	tempArenaInfoDir  string = "temp_arena_info"
	ReplayDir         string = "replays"
	TempArenaInfoName string = "tempArenaInfo.json"
)

type TempArenaInfo struct{}

func NewTempArenaInfo() *TempArenaInfo {
	return &TempArenaInfo{}
}

func (t *TempArenaInfo) Get(installPath string) (vo.TempArenaInfo, error) {
	errDetail := apperr.Tai.Get
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
		return tempArenaInfo, errors.WithStack(errDetail.WithRaw(err))
	}

	tempArenaInfo, err = decideTempArenaInfo(tempArenaInfoPaths)
	if err != nil {
		return tempArenaInfo, err
	}

	return tempArenaInfo, nil
}

func (t *TempArenaInfo) Save(tempArenaInfo vo.TempArenaInfo) error {
	errDetail := apperr.Tai.Save

	_ = os.Mkdir(tempArenaInfoDir, 0o755)

	date, err := time.Parse("2006-01-02 15:04:05", tempArenaInfo.FormattedDateTime())
	if err != nil {
		return errors.WithStack(errDetail.WithRaw(err))
	}
	file, err := os.Create(filepath.Join(tempArenaInfoDir, "tempArenaInfo_"+strconv.FormatInt(date.Unix(), 10)+".json"))
	if err != nil {
		return errors.WithStack(errDetail.WithRaw(err))
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(tempArenaInfo); err != nil {
		return errors.WithStack(errDetail.WithRaw(err))
	}

	return nil
}

//nolint:cyclop
func decideTempArenaInfo(paths []string) (vo.TempArenaInfo, error) {
	errDetail := apperr.Tai.Get
	var tempArenaInfo vo.TempArenaInfo

	size := len(paths)

	if size == 0 {
		return tempArenaInfo, errors.WithStack(errDetail.WithRaw(apperr.ErrNoTempArenaInfo))
	}

	if size == 1 {
		data, err := os.ReadFile(paths[0])
		if err != nil {
			return tempArenaInfo, errors.WithStack(errDetail.WithRaw(err))
		}

		if err := json.Unmarshal(data, &tempArenaInfo); err != nil {
			return tempArenaInfo, errors.WithStack(errDetail.WithRaw(err))
		}

		return tempArenaInfo, nil
	}

	var latest *vo.TempArenaInfo
	var latestDate time.Time
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var tempArenaInfo vo.TempArenaInfo
		if err := json.Unmarshal(data, &tempArenaInfo); err != nil {
			continue
		}

		date, err := time.Parse("2006-01-02 15:04:05", tempArenaInfo.FormattedDateTime())
		if err != nil {
			continue
		}

		if date.After(latestDate) {
			latest = &tempArenaInfo
			latestDate = date
		}
	}

	if latest == nil {
		return vo.TempArenaInfo{}, errors.WithStack(errDetail.WithRaw(apperr.ErrNoTempArenaInfo))
	}

	return *latest, nil
}
