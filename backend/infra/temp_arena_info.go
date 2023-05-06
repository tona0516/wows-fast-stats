package infra

import (
	"changeme/backend/apperr"
	"changeme/backend/vo"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/morikuni/failure"
)

const (
	tempArenaInfoDir  string = "temp_arena_info"
	ReplayDir         string = "replays"
	TempArenaInfoName string = "tempArenaInfo.json"
)

type TempArenaInfo struct{}

func (t *TempArenaInfo) Get(installPath string) (vo.TempArenaInfo, error) {
	errCode := apperr.TempArenaInfoGet

	var tempArenaInfo vo.TempArenaInfo
	data, err := os.ReadFile(filepath.Join(installPath, ReplayDir, TempArenaInfoName))
	if err != nil {
		return tempArenaInfo, failure.Translate(err, errCode)
	}

	if err := json.Unmarshal(data, &tempArenaInfo); err != nil {
		return tempArenaInfo, failure.Translate(err, errCode)
	}

	return tempArenaInfo, nil
}

func (t *TempArenaInfo) Save(tempArenaInfo vo.TempArenaInfo) error {
	errCode := apperr.TempArenaInfoSave

	_ = os.Mkdir(tempArenaInfoDir, 0o755)

	date, err := time.Parse("2006-01-02 15:04:05", tempArenaInfo.FormattedDateTime())
	if err != nil {
		return failure.Translate(err, errCode)
	}
	file, err := os.Create(filepath.Join(tempArenaInfoDir, "tempArenaInfo_"+strconv.FormatInt(date.Unix(), 10)+".json"))
	if err != nil {
		return failure.Translate(err, errCode)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(tempArenaInfo); err != nil {
		return failure.Translate(err, errCode)
	}

	return nil
}
