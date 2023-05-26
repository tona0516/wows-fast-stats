package infra

import (
	"changeme/backend/apperr"
	"changeme/backend/vo"
	"encoding/json"
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
	data, err := os.ReadFile(filepath.Join(installPath, ReplayDir, TempArenaInfoName))
	if err != nil {
		return tempArenaInfo, errors.WithStack(errDetail.WithRaw(err))
	}

	if err := json.Unmarshal(data, &tempArenaInfo); err != nil {
		return tempArenaInfo, errors.WithStack(errDetail.WithRaw(err))
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
