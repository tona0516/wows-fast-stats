package infra

import (
	"changeme/backend/vo"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type TempArenaInfo struct{}

func (t *TempArenaInfo) Get(installPath string) (vo.TempArenaInfo, error) {
	var tempArenaInfo vo.TempArenaInfo
	data, err := os.ReadFile(filepath.Join(installPath, "replays", "tempArenaInfo.json"))
	if err != nil {
		return tempArenaInfo, err
	}

	err = json.Unmarshal(data, &tempArenaInfo)
	if err != nil {
		return tempArenaInfo, err
	}

	return tempArenaInfo, nil
}

func (t *TempArenaInfo) Save(tempArenaInfo vo.TempArenaInfo) error {
    os.Mkdir("temp_arena_info", 0755)

    date, err := time.Parse("2006-01-02 15:04:05", tempArenaInfo.FormattedDateTime())
    if err != nil {
        return err
    }
    file, err := os.Create(filepath.Join("temp_arena_info", "tempArenaInfo_" + strconv.FormatInt(date.Unix(), 10) + ".json"))
	if err != nil {
		return err
	}
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    err = encoder.Encode(tempArenaInfo)
    return err
}
