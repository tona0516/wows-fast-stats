package usecase

import (
	"context"
	"fmt"
	"sync"
	"wfs/backend/adapter"
	"wfs/backend/core"

	"github.com/samber/do/v2"
	"golang.org/x/sync/errgroup"
)

type Prefetch struct {
	cacheStore      adapter.CacheStore
	wargamingClient adapter.WargamingClient
	numbersClient   adapter.NumbersClient
	logger          adapter.Logger
}

func NewPrefetch(i do.Injector) (*Prefetch, error) {
	return &Prefetch{
		cacheStore:      do.MustInvoke[adapter.CacheStore](i),
		wargamingClient: do.MustInvoke[adapter.WargamingClient](i),
		numbersClient:   do.MustInvoke[adapter.NumbersClient](i),
		logger:          do.MustInvoke[adapter.Logger](i),
	}, nil
}

func (p *Prefetch) Invoke(ctx context.Context) (*core.PrefetchResult, error) {
	eg, egCtx := errgroup.WithContext(ctx)

	var warships core.Warships
	eg.Go(func() error {
		var err error
		elapsed := measure(func() {
			warships, err = p.fetchWarships(egCtx)
		})
		p.logger.Debug(fmt.Sprintf("fetchWarships took %dms", elapsed), nil)
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return &core.PrefetchResult{
		Warships: warships,
	}, nil
}

func (p *Prefetch) fetchWarships(ctx context.Context) (core.Warships, error) {
	eg, egCtx := errgroup.WithContext(ctx)

	var encycShips map[int]core.WGEncycShips
	eg.Go(func() error {
		resp, err := p.fetchEncycShips(egCtx)
		if err != nil {
			return err
		}
		encycShips = resp
		return nil
	})

	var expectedStats core.NSExpectedStats
	eg.Go(func() error {
		resp, err := p.numbersClient.ExpectedStats(egCtx)
		if err != nil {
			return err
		}
		expectedStats = resp
		return nil
	})

	if err := eg.Wait(); err != nil {
		cache, errCache := p.cacheStore.Warships()
		if errCache != nil {
			return nil, err
		}

		return cache, nil
	}

	warships := p.composeWarships(encycShips, expectedStats)
	p.cacheStore.SetWarships(warships)

	return warships, nil
}

func (p *Prefetch) fetchEncycShips(ctx context.Context) (map[int]core.WGEncycShips, error) {
	result := make(map[int]core.WGEncycShips)

	var mu sync.Mutex
	fetch := func(ctx context.Context,
		page int) (int, error) {
		res, err := p.wargamingClient.EncycShips(ctx, page)
		if err != nil {
			return 0, err
		}

		mu.Lock()
		result[page] = res
		mu.Unlock()

		return res.Meta.PageTotal, nil
	}

	pageTotal, err := fetch(ctx, 1)
	if err != nil {
		return nil, err
	}

	eg, egCtx := errgroup.WithContext(ctx)

	for i := 2; i < pageTotal+1; i++ {
		eg.Go(func() error {
			_, err := fetch(egCtx, i)
			if err != nil {
				return err
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *Prefetch) composeWarships(
	encycShips map[int]core.WGEncycShips,
	expectedStats core.NSExpectedStats,
) core.Warships {
	warships := make(core.Warships)
	for _, resp := range encycShips {
		for shipID, ship := range resp.Data {
			var serverAverage *core.ServerAverage

			expected, ok := expectedStats.Data[shipID]
			if ok {
				serverAverage = &core.ServerAverage{
					Damage:  expected.AverageDamageDealt,
					Frags:   expected.AverageFrags,
					WinRate: expected.WinRate,
				}
			}

			warship := core.NewWarship(
				shipID,
				ship.Name,
				ship.Tier,
				core.NewShipType(ship.Type),
				core.Nation(ship.Nation),
				ship.IsPremium,
				serverAverage,
			)
			warships[shipID] = *warship
		}
	}

	return warships
}
