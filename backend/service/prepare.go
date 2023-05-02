package service

import (
	"changeme/backend/infra"
	"changeme/backend/vo"
	"os"
	"sync"
)

type Prepare struct {
    parallels uint
    wargaming infra.Wargaming
    numbers infra.Numbers
}

func NewPrepare(parallels uint, wargaming infra.Wargaming, numbers infra.Numbers) *Prepare {
    return &Prepare{
        parallels: parallels,
        wargaming: wargaming,
        numbers: numbers,
    }
}

func (p *Prepare) FetchCachable() error {
    if err := p.deleteOldCache(); err != nil {
        return err
    }

    warshipResult := make(chan error)
	expectedStatsResult := make(chan error)
    battleArenasResult := make(chan error)
    battleTypesResult := make(chan error)

    go p.warship(warshipResult)
    go p.expectedStats(expectedStatsResult)
    go p.battleArenas(battleArenasResult)
    go p.battleTypes(battleTypesResult)

    if err := <-warshipResult; err != nil {
        return err
    }

    if err := <-expectedStatsResult; err != nil {
        return err
    }

    if err := <-battleArenasResult; err != nil {
        return err
    }

    if err := <-battleTypesResult; err != nil {
        return err
    }

    return nil
}

func (p *Prepare) deleteOldCache() error {
    return os.RemoveAll(infra.CACHE_DIRECTORY)
}

func (p *Prepare) warship(result chan error) {
	warships := make(map[int]vo.Warship, 0)

	res, err := p.wargaming.EncyclopediaShips(1)
	if err != nil {
		result <- err
		return
	}
	pageTotal := res.Meta.PageTotal

	var mu sync.Mutex
	limit := make(chan struct{}, p.parallels)
	wg := sync.WaitGroup{}
	for i := 1; i <= pageTotal; i++ {
		limit <- struct{}{}
		wg.Add(1)
		go func(pageNo int) {
			defer func() {
				wg.Done()
				<-limit
			}()

			encyclopediaShips, err := p.wargaming.EncyclopediaShips(pageNo)
			if err != nil {
				result <- err
				return
			}

			for shipID, shipInfo := range encyclopediaShips.Data {
				mu.Lock()
				warships[shipID] = vo.Warship{
					Name:   shipInfo.Name,
					Tier:   shipInfo.Tier,
					Type:   shipInfo.Type,
					Nation: shipInfo.Nation,
				}
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()

    unregistered := infra.Unregistered{}
    unregisteredShipInfo, err := unregistered.Warship()
    if err != nil {
        result <- err
		return
    }
    for k, v := range unregisteredShipInfo {
        if _, ok := warships[k]; !ok {
            warships[k] = v
        }
    }

    cache := infra.Cache[map[int]vo.Warship]{Name: "warship"}
    result <- cache.Serialize(warships)
}

func (p *Prepare) expectedStats(result chan error) {
	expectedStats, err := p.numbers.ExpectedStats()
	if err != nil {
		result <- err
		return
	}

    cache := infra.Cache[vo.NSExpectedStats]{Name: "expectedstats"}
    result <- cache.Serialize(*expectedStats)
}

func (p *Prepare) battleArenas(result chan error) {
	battleArenas, err := p.wargaming.BattleArenas()
	if err != nil {
		result <- err
		return
	}

    cache := infra.Cache[vo.WGBattleArenas]{Name: "battlearenas"}
	result <- cache.Serialize(battleArenas)
}

func (p *Prepare) battleTypes(result chan error) {
	battleTypes, err := p.wargaming.BattleTypes()
	if err != nil {
		result <- err
		return
	}

    cache := infra.Cache[vo.WGBattleTypes]{Name: "battletypes"}

	result <- cache.Serialize(battleTypes)
}
