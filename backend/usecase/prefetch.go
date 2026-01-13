package usecase

import (
	"context"
	"sync"
	"wfs/backend/adapter"
	"wfs/backend/data"

	"github.com/morikuni/failure"
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

func (p *Prefetch) Invoke(ctx context.Context) (*data.PrefetchResult, error) {
	eg, egCtx := errgroup.WithContext(ctx)

	var warships data.Warships
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
		return nil, failure.Wrap(err)
	}

	return &data.PrefetchResult{
		Warships:     warships,
		BattleArenas: battleArenas,
		BattleTypes:  battleTypes,
	}, nil
}

func (p *Prefetch) fetchWarships(ctx context.Context) (data.Warships, error) {
	eg, egCtx := errgroup.WithContext(ctx)

	var encycShips map[int]data.WGEncycShips
	eg.Go(func() error {
		resp, err := p.fetchEncycShips(egCtx)
		if err != nil {
			return failure.Wrap(err)
		}
		encycShips = resp
		return nil
	})

	var expectedStats data.NSExpectedStats
	eg.Go(func() error {
		resp, err := p.numbersClient.ExpectedStats(egCtx)
		if err != nil {
			return failure.Wrap(err)
		}
		expectedStats = resp
		return nil
	})

	if err := eg.Wait(); err != nil {
		cache, errCache := p.cacheStore.Warships()
		if errCache != nil {
			return nil, failure.Wrap(err)
		}

		return cache, nil
	}

	warships := p.composeWarships(encycShips, expectedStats)
	_ = p.cacheStore.SetWarships(warships)

	return warships, nil
}

func (p *Prefetch) fetchBattleArenas(ctx context.Context) (map[int]string, error) {
	resp, err := p.wargamingClient.BattleArenas(ctx)
	if err != nil {
		cache, errCache := p.cacheStore.BattleArenas()
		if errCache != nil {
			return nil, failure.Wrap(err)
		}
		return cache, nil
	}

	result := make(map[int]string)
	for id, arena := range resp.Data {
		result[id] = arena.Name
	}

	_ = p.cacheStore.SetBattleArenas(result)

	return result, nil
}

func (p *Prefetch) fetchBattleTypes(ctx context.Context) (map[string]string, error) {
	resp, err := p.wargamingClient.BattleTypes(ctx)
	if err != nil {
		cache, errCache := p.cacheStore.BattleTypes()
		if errCache != nil {
			return nil, failure.Wrap(err)
		}
		return cache, nil
	}

	result := make(map[string]string)
	for key, battleType := range resp.Data {
		result[key] = battleType.Name
	}

	_ = p.cacheStore.SetBattleTypes(result)

	return result, nil
}

func (p *Prefetch) fetchEncycShips(ctx context.Context) (map[int]data.WGEncycShips, error) {
	result := make(map[int]data.WGEncycShips)

	var mu sync.Mutex
	fetch := func(ctx context.Context,
		page int) (int, error) {
		res, err := p.wargamingClient.EncycShips(ctx, page)
		if err != nil {
			return 0, failure.Wrap(err)
		}

		mu.Lock()
		result[page] = res
		mu.Unlock()

		return res.Meta.PageTotal, nil
	}

	pageTotal, err := fetch(ctx, 1)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	eg, egCtx := errgroup.WithContext(ctx)

	for i := 2; i < pageTotal+1; i++ {
		eg.Go(func() error {
			_, err := fetch(egCtx, i)
			if err != nil {
				return failure.Wrap(err)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, failure.Wrap(err)
	}

	return result, nil
}

func (p *Prefetch) composeWarships(
	encycShips map[int]data.WGEncycShips,
	expectedStats data.NSExpectedStats,
) data.Warships {
	warships := make(data.Warships)
	for _, resp := range encycShips {
		for shipID, ship := range resp.Data {
			var serverAverageDamage, serverAverageFrags, serverAverageWinRate float64

			expected, ok := expectedStats.Data[shipID]
			if ok {
				serverAverageDamage = expected.AverageDamageDealt
				serverAverageFrags = expected.AverageFrags
				serverAverageWinRate = expected.WinRate
			}

			warship := data.NewWarship(
				shipID,
				ship.Name,
				ship.Tier,
				data.NewShipType(ship.Type),
				data.Nation(ship.Nation),
				ship.IsPremium,
				data.ServerAverage{
					Damage:  serverAverageDamage,
					Frags:   serverAverageFrags,
					WinRate: serverAverageWinRate,
				},
			)
			warships[shipID] = *warship
		}
	}

	return warships
}
