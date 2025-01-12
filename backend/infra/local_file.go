package infra

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strconv"
	"wfs/backend/domain/model"

	"github.com/morikuni/failure"
)

const tempArenaInfoDir string = "temp_arena_info"

type LocalFile struct {
	userConfigPath  string
	alertPlayerPath string
}

func NewLocalFile() *LocalFile {
	return &LocalFile{
		userConfigPath:  filepath.Join(ConfigDir, UserConfigFile),
		alertPlayerPath: filepath.Join(ConfigDir, AlertPlayerFile),
	}
}

func (l *LocalFile) SaveScreenshot(path string, base64Data string) error {
	dir := filepath.Dir(path)
	_ = os.Mkdir(dir, 0o755)

	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return failure.Wrap(err)
	}

	f, err := os.Create(path)
	if err != nil {
		return failure.Wrap(err)
	}
	defer f.Close()

	_, err = f.Write(data)
	return failure.Wrap(err)
}

func (l *LocalFile) SaveTempArenaInfo(tempArenaInfo model.TempArenaInfo) error {
	path := filepath.Join(tempArenaInfoDir, "tempArenaInfo_"+strconv.FormatInt(tempArenaInfo.Unixtime(), 10)+".json")
	return writeJSON(path, tempArenaInfo)
}
