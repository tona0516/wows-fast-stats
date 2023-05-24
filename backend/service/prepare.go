package service

import (
	"changeme/backend/apperr"
	"changeme/backend/infra"
	"changeme/backend/vo"
	"os"
	"sync"

	"github.com/pkg/errors"
)

type Prepare struct {
	parallels    uint
	wargaming    infra.WargamingInterface
	numbers      infra.NumbersInterface
	unregistered infra.UnregisteredInterface
	caches       infra.Caches
}

func NewPrepare(
	parallels uint,
	wargaming infra.WargamingInterface,
	numbers infra.NumbersInterface,
	unregistered infra.UnregisteredInterface,
	caches infra.Caches,
) *Prepare {
	return &Prepare{
		parallels:    parallels,
		wargaming:    wargaming,
		numbers:      numbers,
		unregistered: unregistered,
		caches:       caches,
	}
}

func (p *Prepare) FetchCachable(result chan error) {
	if err := p.deleteOldCache(); err != nil {
		result <- err

		return
	}

	fns := [](func(chan error)){
		p.warship,
		p.expectedStats,
		p.battleArenas,
		p.battleTypes,
	}

	errs := make([](*chan error), 0)
	for _, fn := range fns {
		fn := fn
		err := make(chan error)
		go fn(err)
		errs = append(errs, &err)
	}

	for _, err := range errs {
		if err := <-*err; err != nil {
			result <- err

			return
		}
	}

	result <- nil
}

func (p *Prepare) deleteOldCache() error {
	if err := os.RemoveAll(p.caches.Dir); err != nil {
		return errors.WithStack(apperr.SrvPrep.DeleteCache.WithRaw(err))
	}

	return nil
}

func (p *Prepare) warship(result chan error) {
	warships := make(map[int]vo.Warship, 0)

	res, err := p.wargaming.EncyclopediaShips(1)
	if err != nil {
		result <- err

		return
	}

	var mu sync.Mutex
	pages := makeRange(1, res.Meta.PageTotal)
	err = doParallel(p.parallels, pages, func(page int) error {
		encyclopediaShips, err := p.wargaming.EncyclopediaShips(page)
		if err != nil {
			return err
		}

		for shipID, warship := range encyclopediaShips.Data {
			mu.Lock()
			warships[shipID] = vo.Warship{
				Name:   warship.Name,
				Tier:   warship.Tier,
				Type:   vo.NewShipType(warship.Type),
				Nation: warship.Nation,
			}
			mu.Unlock()
		}

		return nil
	})
	if err != nil {
		result <- err

		return
	}

	unregisteredShipInfo, err := p.unregistered.Warship()
	if err != nil {
		result <- err

		return
	}
	for k, v := range unregisteredShipInfo {
		if _, ok := warships[k]; !ok {
			warships[k] = v
		}
	}

	result <- p.caches.Warship.Serialize(warships)
}

func (p *Prepare) expectedStats(result chan error) {
	expectedStats, err := p.numbers.ExpectedStats()
	if err != nil {
		result <- err

		return
	}

	result <- p.caches.ExpectedStats.Serialize(expectedStats)
}

func (p *Prepare) battleArenas(result chan error) {
	battleArenas, err := p.wargaming.BattleArenas()
	if err != nil {
		result <- err

		return
	}

	result <- p.caches.BattleArenas.Serialize(battleArenas)
}

func (p *Prepare) battleTypes(result chan error) {
	battleTypes, err := p.wargaming.BattleTypes()
	if err != nil {
		result <- err

		return
	}

	result <- p.caches.BattleTypes.Serialize(battleTypes)
}
