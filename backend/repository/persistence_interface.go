package repository

import (
	"wfs/backend/data"
	"wfs/backend/domain"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type PersistenceInterface interface {
	// IGN
	LoadOwnIGN() (string, error)
	SaveOwnIGN(ign string) error
	// Expected Stats
	LoadExpectedStats() (data.ExpectedStats, error)
	SaveExpectedStats(data data.ExpectedStats) error
	// User Config.
	LoadUserConfig() (domain.UserConfig, error)
	SaveUserConfig(data domain.UserConfig) error
	// Blacklist.
	LoadBlackList() (domain.BlackList, error)
	SaveBlackList(data domain.BlackList) error
}
