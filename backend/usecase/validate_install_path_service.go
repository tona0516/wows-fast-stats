package usecase

import (
	"os"
	"path/filepath"
	"wfs/backend/core"

	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
)

type ValidateInstallPathService struct {
	gameClientFile string
}

func NewValidateInstallPathService(i do.Injector) (*ValidateInstallPathService, error) {
	return &ValidateInstallPathService{
		gameClientFile: "WorldOfWarships.exe",
	}, nil
}

func (s *ValidateInstallPathService) Validate(input string) error {
	if input == "" {
		return failure.New(core.ErrEmptyInstallPath)
	}

	if _, err := os.Stat(filepath.Join(input, s.gameClientFile)); err != nil {
		return failure.Translate(err, core.ErrNotGameClientPath)
	}

	return nil
}
