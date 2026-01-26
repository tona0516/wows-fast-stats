package usecase

import (
	"context"
	"fmt"
	"wfs/backend/adapter"
	"wfs/backend/core"
	"wfs/backend/usecase/service"

	"github.com/samber/do/v2"
	"golang.org/x/sync/errgroup"
)

type FetchBattle struct {
	wargamingClient adapter.WargamingClient
	cacheStore      adapter.CacheStore
	statsFetcher    *service.StatsFetcher
	clanFetcher     *service.ClanFetcher
	badgeFetcher    *service.BadgeFetcher
	logger          adapter.Logger
}

func NewFetchBattle(i do.Injector) (*FetchBattle, error) {
	return &FetchBattle{
		wargamingClient: do.MustInvoke[adapter.WargamingClient](i),
		cacheStore:      do.MustInvoke[adapter.CacheStore](i),
		statsFetcher:    do.MustInvoke[*service.StatsFetcher](i),
		clanFetcher:     do.MustInvoke[*service.ClanFetcher](i),
		badgeFetcher:    do.MustInvoke[*service.BadgeFetcher](i),
		logger:          do.MustInvoke[adapter.Logger](i),
	}, nil
}

func (b *FetchBattle) Invoke(
	ctx context.Context,
	tempArenaInfo core.TempArenaInfo,
	prefetchResult *core.PrefetchResult,
) (*core.Battle, error) {
	b.cacheStore.SetOwnIGN(tempArenaInfo.PlayerName)
	b.logger.SetOwnIGN(tempArenaInfo.PlayerName)

	accountNames := tempArenaInfo.AccountNames()

	var accountList core.WGAccountList
	var err error
	elapsed := measure(func() {
		accountList, err = b.wargamingClient.AccountList(ctx, accountNames)
	})
	b.logger.Debug(fmt.Sprintf("AccountList took %dms", elapsed), nil)

	if err != nil {
		return nil, err
	}
	accountIDs := accountList.AccountIDs()

	eg, egCtx := errgroup.WithContext(ctx)

	var accountInfo core.WGAccountInfo
	eg.Go(func() error {
		var err error
		elapsed := measure(func() {
			accountInfo, err = b.wargamingClient.AccountInfo(egCtx, accountIDs)
		})
		b.logger.Debug(fmt.Sprintf("AccountInfo took %dms", elapsed), nil)
		return err
	})

	var allShipStats core.AllPlayerShipStats
	eg.Go(func() error {
		var err error
		elapsed := measure(func() {
			allShipStats, err = b.statsFetcher.FetchAll(egCtx, accountIDs)
		})
		b.logger.Debug(fmt.Sprintf("statsFetcher.FetchAll took %dms", elapsed), nil)
		return err
	})

	var allShipBadges core.AllPlayerShipBadges
	eg.Go(func() error {
		var err error
		elapsed := measure(func() {
			allShipBadges, err = b.badgeFetcher.FetchAll(egCtx, accountIDs)
		})
		b.logger.Debug(fmt.Sprintf("badgeFetcher.FetchAll took %dms", elapsed), nil)
		return err
	})

	var clans core.Clans
	eg.Go(func() error {
		var err error
		elapsed := measure(func() {
			clans, err = b.clanFetcher.FetchAll(egCtx, accountIDs)
		})
		b.logger.Debug(fmt.Sprintf("clanFetcher.FetchAll took %dms", elapsed), nil)
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
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

	return &result, nil
}
