package usecase

import (
	"context"
	"os"
	"path/filepath"
	"wfs/internal/infra"

	"github.com/morikuni/failure"
)

const gameClientFile = "WorldOfWarships.exe"

type openDirectoryDialogFunc func(ctx context.Context) (string, error)

type InstallPathSetting struct {
	localStorage        infra.ConfigStore
	openDirectoryDialog openDirectoryDialogFunc
}

func NewInstallPathSetting(
	localStorage infra.ConfigStore,
	openDirectoryDialogFunc openDirectoryDialogFunc,
) *InstallPathSetting {
	return &InstallPathSetting{
		localStorage:        localStorage,
		openDirectoryDialog: openDirectoryDialogFunc,
	}
}

func (s *InstallPathSetting) Invoke(ctx context.Context) (bool, error) {
	selected, err := s.openDirectoryDialog(ctx)
	if err != nil {
		return false, failure.Wrap(err)
	}

	if selected == "" {
		return false, nil
	}

	if _, err := os.Stat(filepath.Join(selected, gameClientFile)); err != nil {
		return false, failure.Wrap(err)
	}

	config, err := s.localStorage.UserConfig()
	if err != nil {
		return false, failure.Wrap(err)
	}

	config.InstallPath = selected

	if err = s.localStorage.SetUserConfig(config); err != nil {
		return false, failure.Wrap(err)
	}

	return true, nil
}
