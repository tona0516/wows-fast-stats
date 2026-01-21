package usecase

import (
	"wfs/backend/adapter"
	"wfs/backend/core"

	"github.com/samber/do/v2"
)

type SavePref struct {
	prefStore adapter.PrefStore
}

func NewSavePref(i do.Injector) (*SavePref, error) {
	return &SavePref{
		prefStore: do.MustInvoke[adapter.PrefStore](i),
	}, nil
}

func (sp *SavePref) Invoke(config core.Pref) error {
	if err := sp.prefStore.SetPref(config); err != nil {
		return err
	}

	return nil
}
