package usecase

import (
	"context"
	"regexp"
	"wfs/backend/adapter"
	"wfs/backend/core"

	"github.com/samber/do/v2"
	"golang.org/x/sync/errgroup"
)

// URLを検出する正規表現パターン.
var urlRegex = regexp.MustCompile(`https?://[^\s]+`)

type FetchBattle struct {
	wargamingClient adapter.WargamingClient
	wails           adapter.Wails
	cacheStore      adapter.CacheStore
	logger          adapter.Logger
	statsService    *statsService
	clanService     *clanService
	badgeService    *badgeService
}

func NewFetchBattle(i do.Injector) (*FetchBattle, error) {
	return &FetchBattle{
		wargamingClient: do.MustInvoke[adapter.WargamingClient](i),
		wails:           do.MustInvoke[adapter.Wails](i),
		cacheStore:      do.MustInvoke[adapter.CacheStore](i),
		logger:          do.MustInvoke[adapter.Logger](i),
		statsService:    do.MustInvoke[*statsService](i),
		clanService:     do.MustInvoke[*clanService](i),
		badgeService:    do.MustInvoke[*badgeService](i),
	}, nil
}

func (b *FetchBattle) Invoke(
	ctx context.Context,
	tempArenaInfo core.TempArenaInfo,
	prefetchResult *core.PrefetchResult,
) {
	b.cacheStore.SetOwnIGN(tempArenaInfo.PlayerName)
	b.logger.SetOwnIGN(tempArenaInfo.PlayerName)

	accountNames := tempArenaInfo.AccountNames()

	var accountList core.WGAccountList
	var err error
	measure("AccountList", func() {
		accountList, err = b.wargamingClient.AccountList(ctx, accountNames)
	})
	if err != nil {
		b.wails.EmitEvent(ctx, EventErr, err)
		return
	}
	accountIDs := accountList.AccountIDs()

	eg, egCtx := errgroup.WithContext(ctx)

	var accountInfo core.WGAccountInfo
	eg.Go(func() error {
		var err error
		measure("AccountInfo", func() {
			accountInfo, err = b.wargamingClient.AccountInfo(egCtx, accountIDs)
		})
		return err
	})

	var allShipStats core.AllPlayerShipStats
	eg.Go(func() error {
		var err error
		measure("statsService.fetchAll", func() {
			allShipStats, err = b.statsService.fetchAll(egCtx, accountIDs)
		})
		return err
	})

	var allShipBadges core.AllPlayerShipBadges
	eg.Go(func() error {
		var err error
		measure("badgeService.fetchAll", func() {
			allShipBadges, err = b.badgeService.fetchAll(egCtx, accountIDs)
		})
		return err
	})

	var clans core.Clans
	eg.Go(func() error {
		var err error
		measure("clanService.fetchAll", func() {
			clans, err = b.clanService.fetchAll(egCtx, accountIDs)
		})
		return err
	})

	if err := eg.Wait(); err != nil {
		b.wails.EmitEvent(ctx, EventErr, err)
		return
	}

	result := core.NewBattle(
		prefetchResult,
		tempArenaInfo,
		accountInfo,
		accountList,
		clans,
		allShipStats,
		allShipBadges,
	)

	b.wails.EmitEvent(ctx, EventFetchDone, result)
}
