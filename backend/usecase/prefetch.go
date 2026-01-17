package usecase

import (
	"context"
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
}

func NewPrefetch(i do.Injector) (*Prefetch, error) {
	return &Prefetch{
		cacheStore:      do.MustInvoke[adapter.CacheStore](i),
		wargamingClient: do.MustInvoke[adapter.WargamingClient](i),
		numbersClient:   do.MustInvoke[adapter.NumbersClient](i),
	}, nil
}

func (p *Prefetch) Invoke(ctx context.Context) (*core.PrefetchResult, error) {
	eg, egCtx := errgroup.WithContext(ctx)

	var warships core.Warships
	eg.Go(func() error {
		var err error
		measure("fetchWarships", func() {
			warships, err = p.fetchWarships(egCtx)
		})
		return err
	})

	var battleArenas map[int]string
	eg.Go(func() error {
		var err error
		measure("fetchBattleArenas", func() {
			battleArenas, err = p.fetchBattleArenas(egCtx)
		})
		return err
	})

	var battleTypes map[string]string
	eg.Go(func() error {
		var err error
		measure("fetchBattleTypes", func() {
			battleTypes, err = p.fetchBattleTypes(egCtx)
		})
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return &core.PrefetchResult{
		Warships:     warships,
		BattleArenas: battleArenas,
		BattleTypes:  battleTypes,
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

func (p *Prefetch) fetchBattleArenas(ctx context.Context) (map[int]string, error) {
	resp, err := p.wargamingClient.BattleArenas(ctx)
	if err != nil {
		cache, errCache := p.cacheStore.BattleArenas()
		if errCache != nil {
			return nil, err
		}
		return cache, nil
	}

	result := make(map[int]string)
	for id, arena := range resp.Data {
		result[id] = arena.Name
	}

	p.cacheStore.SetBattleArenas(result)

	return result, nil
}

func (p *Prefetch) fetchBattleTypes(ctx context.Context) (map[string]string, error) {
	resp, err := p.wargamingClient.BattleTypes(ctx)
	if err != nil {
		cache, errCache := p.cacheStore.BattleTypes()
		if errCache != nil {
			return nil, err
		}
		return cache, nil
	}

	result := make(map[string]string)
	for key, battleType := range resp.Data {
		result[key] = battleType.Name
	}

	p.cacheStore.SetBattleTypes(result)

	return result, nil
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
