package service

import (
	"os"
	"path/filepath"
	"wfs/backend/core"

	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
)

type GameClientPathValidator struct {
	gameClientFile string
}

func NewGameClientPathValidator(i do.Injector) (*GameClientPathValidator, error) {
	return &GameClientPathValidator{
		gameClientFile: "WorldOfWarships.exe",
	}, nil
}

func (v *GameClientPathValidator) Validate(input string) error {
	if input == "" {
		return failure.New(core.ErrEmptyGameClientPath)
	}

	if _, err := os.Stat(filepath.Join(input, v.gameClientFile)); err != nil {
		return failure.Translate(err, core.ErrNotGameClientPath)
	}

	return nil
}
