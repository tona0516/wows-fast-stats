package adapter

import (
	"context"
	"wfs/backend/data"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type Wails interface {
	EmitEvent(ctx context.Context, eventName string, optionalData ...any)
	OpenDirectoryDialog(ctx context.Context) (string, error)
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type CacheStore interface {
	OwnIGN() (string, error)
	SetOwnIGN(ign string)

	Warships() (data.Warships, error)
	SetWarships(data data.Warships)

	BattleArenas() (map[int]string, error)
	SetBattleArenas(data map[int]string)

	BattleTypes() (map[string]string, error)
	SetBattleTypes(data map[string]string)
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type PrefStore interface {
	Pref() (data.Pref, error)
	SetPref(data data.Pref) error
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type ReplayReader interface {
	TempArenaInfo(installPath string) (data.TempArenaInfo, error)
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type ClanClient interface {
	ClanAutoComplete(ctx context.Context, search string) (data.ClanAutocomplete, error)
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type DiscordClient interface {
	Comment(ctx context.Context, message string) error
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type GithubClient interface {
	LatestRelease(ctx context.Context) (data.GHLatestRelease, error)
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type NumbersClient interface {
	ExpectedStats(ctx context.Context) (data.NSExpectedStats, error)
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type WargamingClient interface {
	AccountInfo(ctx context.Context, accountIDs []data.AccountID) (data.WGAccountInfo, error)
	AccountList(ctx context.Context, accountNames []string) (data.WGAccountList, error)
	ClansAccountInfo(ctx context.Context, accountIDs []data.AccountID) (data.WGClansAccountInfo, error)
	ClansInfo(ctx context.Context, clanIDs []data.ClanID) (data.WGClansInfo, error)
	EncycShips(ctx context.Context, pageNo int) (data.WGEncycShips, error)
	ShipsStats(ctx context.Context, accountID data.AccountID) (data.WGShipsStats, error)
	BattleArenas(ctx context.Context) (data.WGBattleArenas, error)
	BattleTypes(ctx context.Context) (data.WGBattleTypes, error)
	ShipsBadges(ctx context.Context, accountID data.AccountID) (data.WGShipsBadges, error)
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type Logger interface {
	SetOwnIGN(ownIGN string)
	Debug(message string, contexts map[string]string)
	Info(message string, contexts map[string]string)
	Error(err error, contexts map[string]string)
}
