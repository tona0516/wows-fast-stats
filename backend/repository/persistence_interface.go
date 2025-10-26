package repository

import (
	"wfs/backend/data"
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
	LoadUserConfig() (data.UserConfig, error)
	SaveUserConfig(data data.UserConfig) error
	// Blacklist.
	LoadBlackList() (data.BlackList, error)
	SaveBlackList(data data.BlackList) error
}
