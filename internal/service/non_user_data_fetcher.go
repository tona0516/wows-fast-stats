package service

import (
	"sync"
	"wfs/internal/adapter"
	"wfs/internal/data"

	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
	"golang.org/x/sync/errgroup"
)

type NonUserData struct {
	Warships     data.Warships
	BattleArenas map[int]string
	BattleTypes  map[string]string
}

type NonUserDataFetcher struct {
	cacheStore      adapter.CacheStore
	wargamingClient adapter.WargamingClient
	numbersClient   adapter.NumbersClient
}

func NewNonUserDataFetcher(i do.Injector) (*NonUserDataFetcher, error) {
	return &NonUserDataFetcher{
		cacheStore:      do.MustInvoke[adapter.CacheStore](i),
		wargamingClient: do.MustInvoke[adapter.WargamingClient](i),
		numbersClient:   do.MustInvoke[adapter.NumbersClient](i),
	}, nil
}

func (f *NonUserDataFetcher) Fetch() (*NonUserData, error) {
	eg := errgroup.Group{}

	var warships data.Warships
	eg.Go(func() error {
		var err error
		warships, err = f.fetchWarships()
		return err
	})

	var battleArenas map[int]string
	eg.Go(func() error {
		var err error
		battleArenas, err = f.fetchBattleArenas()
		return err
	})

	var battleTypes map[string]string
	eg.Go(func() error {
		var err error
		battleTypes, err = f.fetchBattleTypes()
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, failure.Wrap(err)
	}

	return &NonUserData{
		Warships:     warships,
		BattleArenas: battleArenas,
		BattleTypes:  battleTypes,
	}, nil
}

func (f *NonUserDataFetcher) fetchWarships() (data.Warships, error) {
	eg := errgroup.Group{}

	var encycShips map[int]data.WGEncycShips
	eg.Go(func() error {
		resp, err := f.fetchEncycShips()
		if err != nil {
			return failure.Wrap(err)
		}
		encycShips = resp
		return nil
	})

	var expectedStats data.NSExpectedStats
	eg.Go(func() error {
		resp, err := f.numbersClient.ExpectedStats()
		if err != nil {
			return failure.Wrap(err)
		}
		expectedStats = resp
		return nil
	})

	if err := eg.Wait(); err != nil {
		cache, errCache := f.cacheStore.Warships()
		if errCache != nil {
			return nil, failure.Wrap(err)
		}

		return cache, nil
	}

	warships := f.composeWarships(encycShips, expectedStats)
	_ = f.cacheStore.SetWarships(warships)

	return warships, nil
}

func (f *NonUserDataFetcher) fetchBattleArenas() (map[int]string, error) {
	resp, err := f.wargamingClient.BattleArenas()
	if err != nil {
		cache, errCache := f.cacheStore.BattleArenas()
		if errCache != nil {
			return nil, failure.Wrap(err)
		}
		return cache, nil
	}

	result := make(map[int]string)
	for id, arena := range resp.Data {
		result[id] = arena.Name
	}

	_ = f.cacheStore.SetBattleArenas(result)

	return result, nil
}

func (f *NonUserDataFetcher) fetchBattleTypes() (map[string]string, error) {
	resp, err := f.wargamingClient.BattleTypes()
	if err != nil {
		cache, errCache := f.cacheStore.BattleTypes()
		if errCache != nil {
			return nil, failure.Wrap(err)
		}
		return cache, nil
	}

	result := make(map[string]string)
	for key, battleType := range resp.Data {
		result[key] = battleType.Name
	}

	_ = f.cacheStore.SetBattleTypes(result)

	return result, nil
}

func (f *NonUserDataFetcher) fetchEncycShips() (map[int]data.WGEncycShips, error) {
	result := make(map[int]data.WGEncycShips)

	var mu sync.Mutex
	fetch := func(page int) (int, error) {
		res, err := f.wargamingClient.EncycShips(page)
		if err != nil {
			return 0, failure.Wrap(err)
		}

		mu.Lock()
		result[page] = res
		mu.Unlock()

		return res.Meta.PageTotal, nil
	}

	pageTotal, err := fetch(1)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	eg := errgroup.Group{}
	for i := 2; i < pageTotal+1; i++ {
		eg.Go(func() error {
			_, err := fetch(i)
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

func (f *NonUserDataFetcher) composeWarships(
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
