package infra

import (
	"path/filepath"
	"wfs/backend/config"
	"wfs/backend/core"

	"github.com/samber/do/v2"
)

type PrefStore struct {
	dir      string
	prefFile string
}

func NewPrefStore(i do.Injector) (*PrefStore, error) {
	config := do.MustInvoke[config.Config](i)
	return &PrefStore{
		dir:      config.LocalFile.RootDir,
		prefFile: "pref.json",
	}, nil
}

func (ps *PrefStore) Pref() (core.Pref, error) {
	return readJSON[core.Pref](filepath.Join(ps.dir, ps.prefFile))
}

func (ps *PrefStore) SetPref(data core.Pref) error {
	return writeJSON(filepath.Join(ps.dir, ps.prefFile), data)
}
