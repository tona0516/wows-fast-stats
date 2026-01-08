package gateway

import "wfs/internal/data"

// store and reader interfaces

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type CacheStore interface {
	OwnIGN() (string, error)
	SetOwnIGN(ign string) error

	Warships() (data.Warships, error)
	SetWarships(data data.Warships) error

	BattleArenas() (map[int]string, error)
	SetBattleArenas(data map[int]string) error

	BattleTypes() (map[string]string, error)
	SetBattleTypes(data map[string]string) error
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type ConfigStore interface {
	UserConfig() (data.UserConfig, error)
	SetUserConfig(data data.UserConfig) error
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type ReplayReader interface {
	TempArenaInfo(installPath string) (data.TempArenaInfo, error)
}

// api client interfaces

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type ClanClient interface {
	ClanAutoComplete(search string) (data.ClanAutocomplete, error)
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type DiscordClient interface {
	Comment(message string) error
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type GithubClient interface {
	LatestRelease() (data.GHLatestRelease, error)
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type NumbersClient interface {
	ExpectedStats() (data.NSExpectedStats, error)
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type WargamingClient interface {
	AccountInfo(accountIDs []int) (data.WGAccountInfo, error)
	AccountList(accountNames []string) (data.WGAccountList, error)
	ClansAccountInfo(accountIDs []int) (data.WGClansAccountInfo, error)
	ClansInfo(clanIDs []int) (data.WGClansInfo, error)
	EncycShips(pageNo int) (data.WGEncycShips, error)
	ShipsStats(accountID int) (data.WGShipsStats, error)
	BattleArenas() (data.WGBattleArenas, error)
	BattleTypes() (data.WGBattleTypes, error)
	ShipsBadges(accountID int) (data.WGShipsBadges, error)
}

// logger interface

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type Logger interface {
	SetOwnIGN(ownIGN string)
	Debug(message string, contexts map[string]string)
	Info(message string, contexts map[string]string)
	Error(err error, contexts map[string]string)
}
