package infra

import (
	"path/filepath"
	"wfs/backend/config"
	"wfs/backend/data"

	"github.com/samber/do/v2"
)

type PrefStore struct {
	dir      string
	prefFile string
}

func NewPrefStore(i do.Injector) (*PrefStore, error) {
	config := do.MustInvoke[config.Config](i)
	return &PrefStore{
		dir:      config.LocalFile.ConfigDir,
		prefFile: "pref.json",
	}, nil
}

func (ps *PrefStore) Pref() (data.Pref, error) {
	return readJSON[data.Pref](filepath.Join(ps.dir, ps.prefFile))
}

func (ps *PrefStore) SetPref(data data.Pref) error {
	return writeJSON(filepath.Join(ps.dir, ps.prefFile), data)
}
