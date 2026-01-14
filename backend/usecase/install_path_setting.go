package usecase

import (
	"context"
	"os"
	"path/filepath"
	"wfs/backend/adapter"
	"wfs/backend/data"

	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
)

type InstallPathSetting struct {
	wails          adapter.Wails
	configStore    adapter.PrefStore
	gameClientFile string
}

func NewInstallPathSetting(i do.Injector) (*InstallPathSetting, error) {
	return &InstallPathSetting{
		wails:          do.MustInvoke[adapter.Wails](i),
		configStore:    do.MustInvoke[adapter.PrefStore](i),
		gameClientFile: "WorldOfWarships.exe",
	}, nil
}

func (s *InstallPathSetting) Invoke(ctx context.Context) (bool, error) {
	selected, err := s.wails.OpenDirectoryDialog(ctx)
	if err != nil {
		return false, err
	}

	if selected == "" {
		return false, nil
	}

	if _, err := os.Stat(filepath.Join(selected, s.gameClientFile)); err != nil {
		return false, failure.Translate(err, data.ErrInvalidInstallPath)
	}

	config, err := s.configStore.Pref()
	if err != nil {
		return false, err
	}

	config.InstallPath = selected

	if err = s.configStore.SetPref(config); err != nil {
		return false, err
	}

	return true, nil
}
