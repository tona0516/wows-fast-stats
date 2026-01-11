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
	configStore adapter.ConfigStore
}

func NewLoadPref(i do.Injector) (*LoadPref, error) {
	return &LoadPref{
		configStore: do.MustInvoke[adapter.ConfigStore](i),
	}, nil
}

func (lp *LoadPref) Invoke() (data.UserConfig, error) {
	config, err := lp.configStore.UserConfig()
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return data.DefaultUserConfig(), nil
		}
		return data.UserConfig{}, failure.Wrap(err)
	}

	return config, nil
}
