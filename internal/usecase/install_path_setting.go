package usecase

import (
	"context"
	"os"
	"path/filepath"
	"wfs/internal/adapter"

	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
)

type InstallPathSetting struct {
	wails          adapter.Wails
	configStore    adapter.ConfigStore
	gameClientFile string
}

func NewInstallPathSetting(i do.Injector) (*InstallPathSetting, error) {
	return &InstallPathSetting{
		wails:          do.MustInvoke[adapter.Wails](i),
		configStore:    do.MustInvoke[adapter.ConfigStore](i),
		gameClientFile: "WorldOfWarships.exe",
	}, nil
}

func (s *InstallPathSetting) Invoke(ctx context.Context) (bool, error) {
	selected, err := s.wails.OpenDirectoryDialog(ctx)
	if err != nil {
		return false, failure.Wrap(err)
	}

	if selected == "" {
		return false, nil
	}

	if _, err := os.Stat(filepath.Join(selected, s.gameClientFile)); err != nil {
		return false, failure.Wrap(err)
	}

	config, err := s.configStore.UserConfig()
	if err != nil {
		return false, failure.Wrap(err)
	}

	config.InstallPath = selected

	if err = s.configStore.SetUserConfig(config); err != nil {
		return false, failure.Wrap(err)
	}

	return true, nil
}
