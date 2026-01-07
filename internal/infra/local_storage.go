package infra

import (
	"io/fs"
	"os"
	"path/filepath"
	"wfs/internal/data"

	"github.com/morikuni/failure"
)

const (
	replaysDir        string = "replays"
	tempArenaInfoFile string = "tempArenaInfo.json"
	userConfigFile    string = "user_config.json"
	ownIGNFile        string = "own_ign.txt"
	warshipsFile      string = "warships.json"
	battleArenasFile  string = "battle_arenas.json"
	battleTypesFile   string = "battle_types.json"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type LocalStorage interface {
	// read
	TempArenaInfo(installPath string) (data.TempArenaInfo, error)
	UserConfig() (data.UserConfig, error)
	OwnIGN() (string, error)
	Warships() (data.Warships, error)
	BattleArenas() (map[int]string, error)
	BattleTypes() (map[string]string, error)
	// write
	SetUserConfig(data data.UserConfig) error
	SetOwnIGN(ign string) error
	SetWarships(data data.Warships) error
	SetBattleArenas(data map[int]string) error
	SetBattleTypes(data map[string]string) error
}

type localStorage struct {
	userDataDir string
	cacheDir    string
}

func NewLocalStorage(
	userDataDir string,
	cacheDir string,
) LocalStorage {
	return &localStorage{
		userDataDir: userDataDir,
		cacheDir:    cacheDir,
	}
}

func (s *localStorage) TempArenaInfo(installPath string) (data.TempArenaInfo, error) {
	var tempArenaInfo data.TempArenaInfo

	tempArenaInfoPaths := []string{}
	root := filepath.Join(installPath, replaysDir)
	if _, err := os.Stat(root); err != nil {
		return tempArenaInfo, failure.Wrap(err)
	}

	err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.Name() != tempArenaInfoFile {
			return nil
		}

		tempArenaInfoPaths = append(tempArenaInfoPaths, path)
		return nil
	})
	if err != nil {
		return tempArenaInfo, failure.Wrap(err)
	}

	return s.decideTempArenaInfo(tempArenaInfoPaths)
}

func (s *localStorage) UserConfig() (data.UserConfig, error) {
	return readJSON[data.UserConfig](filepath.Join(s.userDataDir, userConfigFile))
}

func (s *localStorage) OwnIGN() (string, error) {
	return readString(filepath.Join(s.cacheDir, ownIGNFile))
}

func (s *localStorage) Warships() (data.Warships, error) {
	return readJSON[data.Warships](filepath.Join(s.cacheDir, warshipsFile))
}

func (s *localStorage) BattleArenas() (map[int]string, error) {
	return readJSON[map[int]string](filepath.Join(s.cacheDir, battleArenasFile))
}

func (s *localStorage) BattleTypes() (map[string]string, error) {
	return readJSON[map[string]string](filepath.Join(s.cacheDir, battleTypesFile))
}

func (s *localStorage) SetUserConfig(data data.UserConfig) error {
	return writeJSON(filepath.Join(s.userDataDir, userConfigFile), data)
}

func (s *localStorage) SetOwnIGN(ign string) error {
	return writeString(filepath.Join(s.cacheDir, ownIGNFile), ign)
}

func (s *localStorage) SetWarships(data data.Warships) error {
	return writeJSON(filepath.Join(s.cacheDir, warshipsFile), data)
}

func (s *localStorage) SetBattleArenas(data map[int]string) error {
	return writeJSON(filepath.Join(s.cacheDir, battleArenasFile), data)
}

func (s *localStorage) SetBattleTypes(data map[string]string) error {
	return writeJSON(filepath.Join(s.cacheDir, battleTypesFile), data)
}

func (s *localStorage) decideTempArenaInfo(paths []string) (data.TempArenaInfo, error) {
	var result data.TempArenaInfo
	size := len(paths)

	if size == 0 {
		return result, failure.Wrap(fs.ErrNotExist)
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
		return result, failure.Wrap(fs.ErrNotExist)
	}

	return latest, nil
}
