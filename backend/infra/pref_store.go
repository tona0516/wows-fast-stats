package infra

import (
	"path/filepath"
	"wfs/backend/config"

	"github.com/samber/do/v2"
)

type PrefStore struct {
	dir                string
	gameClientPathFile string
	displayPrefFile    string
}

func NewPrefStore(i do.Injector) (*PrefStore, error) {
	config := do.MustInvoke[config.Config](i)
	return &PrefStore{
		dir:                config.LocalFile.RootDir,
		gameClientPathFile: "game_client_path.txt",
		displayPrefFile:    "display_pref.json",
	}, nil
}

func (ps *PrefStore) GameClientPath() (string, error) {
	return readString(filepath.Join(ps.dir, ps.gameClientPathFile))
}

func (ps *PrefStore) SetGameClientPath(path string) error {
	return writeString(filepath.Join(ps.dir, ps.gameClientPathFile), path)
}

func (ps *PrefStore) DisplayPref() (string, error) {
	return readString(filepath.Join(ps.dir, ps.displayPrefFile))
}

func (ps *PrefStore) SetDisplayPref(pref string) error {
	return writeString(filepath.Join(ps.dir, ps.displayPrefFile), pref)
}
