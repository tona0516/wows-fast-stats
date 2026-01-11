package usecase

import (
	"errors"
	"io/fs"
	"wfs/backend/adapter"
	"wfs/backend/data"

	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
)

type LoadPref struct {
	configStore adapter.PrefStore
}

func NewLoadPref(i do.Injector) (*LoadPref, error) {
	return &LoadPref{
		configStore: do.MustInvoke[adapter.PrefStore](i),
	}, nil
}

func (lp *LoadPref) Invoke() (data.Pref, error) {
	config, err := lp.configStore.Pref()
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return data.DefaultPref(), nil
		}
		return data.Pref{}, failure.Wrap(err)
	}

	return config, nil
}
