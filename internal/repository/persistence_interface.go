package repository

import (
	"wfs/internal/data"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type PersistenceInterface interface {
	// IGN
	LoadOwnIGN() (string, error)
	SaveOwnIGN(ign string) error
	// Expected Stats
	LoadExpectedStats() (data.ExpectedStats, error)
	SaveExpectedStats(data data.ExpectedStats) error
	// User Config.
	LoadUserConfig() (data.UserConfig, error)
	SaveUserConfig(data data.UserConfig) error
}
