package adapter

import (
	"context"
	"wfs/backend/core"
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

	Warships() (core.Warships, error)
	SetWarships(data core.Warships)

	BattleArenas() (map[int]string, error)
	SetBattleArenas(data map[int]string)

	BattleTypes() (map[string]string, error)
	SetBattleTypes(data map[string]string)
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type PrefStore interface {
	GameClientPath() (string, error)
	SetGameClientPath(path string) error
	DisplayPref() (string, error)
	SetDisplayPref(pref string) error
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type ReplayReader interface {
	TempArenaInfo(gameClientPath string) (core.TempArenaInfo, error)
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type ClanClient interface {
	ClanAutoComplete(ctx context.Context, search string) (core.ClanAutocomplete, error)
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type DiscordClient interface {
	Comment(ctx context.Context, message string) error
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type GithubClient interface {
	LatestRelease(ctx context.Context) (core.GHLatestRelease, error)
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type NumbersClient interface {
	ExpectedStats(ctx context.Context) (core.NSExpectedStats, error)
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type WargamingClient interface {
	AccountInfo(ctx context.Context, accountIDs []core.AccountID) (core.WGAccountInfo, error)
	AccountList(ctx context.Context, accountNames []string) (core.WGAccountList, error)
	ClansAccountInfo(ctx context.Context, accountIDs []core.AccountID) (core.WGClansAccountInfo, error)
	ClansInfo(ctx context.Context, clanIDs []core.ClanID) (core.WGClansInfo, error)
	EncycShips(ctx context.Context, pageNo int) (core.WGEncycShips, error)
	ShipsStats(ctx context.Context, accountID core.AccountID) (core.WGShipsStats, error)
	BattleArenas(ctx context.Context) (core.WGBattleArenas, error)
	BattleTypes(ctx context.Context) (core.WGBattleTypes, error)
	ShipsBadges(ctx context.Context, accountID core.AccountID) (core.WGShipsBadges, error)
}

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type Logger interface {
	SetOwnIGN(ownIGN string)
	Debug(message string, contexts map[string]string)
	Info(message string, contexts map[string]string)
	Error(err error, contexts map[string]string)
}
