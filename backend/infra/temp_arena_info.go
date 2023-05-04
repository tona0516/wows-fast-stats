package infra

import (
	"changeme/backend/vo"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type TempArenaInfo struct{}

func (t *TempArenaInfo) Get(installPath string) (vo.TempArenaInfo, error) {
	var tempArenaInfo vo.TempArenaInfo
	data, err := os.ReadFile(filepath.Join(installPath, "replays", "tempArenaInfo.json"))
	if err != nil {
		return tempArenaInfo, errors.WithStack(err)
	}

	err = json.Unmarshal(data, &tempArenaInfo)
	if err != nil {
		return tempArenaInfo, errors.WithStack(err)
	}

	return tempArenaInfo, nil
}

func (t *TempArenaInfo) Save(tempArenaInfo vo.TempArenaInfo) error {
    _ = os.Mkdir("temp_arena_info", 0755)

    date, err := time.Parse("2006-01-02 15:04:05", tempArenaInfo.FormattedDateTime())
    if err != nil {
        return errors.WithStack(err)
    }
    file, err := os.Create(filepath.Join("temp_arena_info", "tempArenaInfo_" + strconv.FormatInt(date.Unix(), 10) + ".json"))
	if err != nil {
		return errors.WithStack(err)
	}
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    err = encoder.Encode(tempArenaInfo)
    return errors.WithStack(err)
}
