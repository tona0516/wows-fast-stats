package infra

import (
	"errors"
	"sync"
	"wfs/backend/apperr"
	"wfs/backend/domain/model"
	"wfs/backend/infra/response"
	"wfs/backend/infra/webapi"

	"github.com/dgraph-io/badger/v4"
	"github.com/morikuni/failure"
)

type WarshipFetcher struct {
	db               *badger.DB
	wargaming        webapi.Wargaming
	numbers          webapi.Numbers
	localDataKeyName string
}

func NewWarshipFetcher(
	db *badger.DB,
	wargaming webapi.Wargaming,
	numbers webapi.Numbers,
) *WarshipFetcher {
	return &WarshipFetcher{
		db:               db,
		wargaming:        wargaming,
		numbers:          numbers,
		localDataKeyName: "warships",
	}
}

func (f *WarshipFetcher) Fetch() (model.Warships, error) {
	cache, errCache := f.readCache()

	currentGameVersion, err := f.wargaming.GameVersion()
	if err != nil {
		return f.toError(cache, errCache, err)
	}

	if currentGameVersion == cache.gameVersion {
		return cache.warships, nil
	}

	encycShipsChan := make(chan model.Result[model.Warships])
	expectedStatsChan := make(chan model.Result[response.ExpectedStats])

	go f.encycShips(encycShipsChan)
	go f.expectedStats(expectedStatsChan)

	ships := <-encycShipsChan
	err = errors.Join(err, ships.Error)

	expectedStats := <-expectedStatsChan
	err = errors.Join(err, expectedStats.Error)
	if err != nil {
		return f.toError(cache, errCache, err)
	}

	warships := ships.Value
	for shipID, ship := range expectedStats.Value {
		if _, ok := warships[shipID]; !ok {
			continue
		}

		w := warships[shipID]
		warships[shipID] = model.Warship{
			ID:            w.ID,
			Name:          w.Name,
			Tier:          w.Tier,
			Type:          w.Type,
			Nation:        w.Nation,
			IsPremium:     w.IsPremium,
			AverageDamage: ship.AverageDamageDealt,
			AverageFrags:  ship.AverageFrags,
			WinRate:       ship.WinRate,
		}
	}

	_ = f.saveCache(warships, currentGameVersion)

	return warships, nil
}

func (f *WarshipFetcher) encycShips(channel chan model.Result[model.Warships]) {
	warships := make(model.Warships)

	var mu sync.Mutex
	fetch := func(page int) (int, error) {
		res, pageTotal, err := f.wargaming.EncycShips(page)
		if err != nil {
			return 0, err
		}

		for shipID, warship := range res.Data {
			mu.Lock()
			warships[shipID] = model.Warship{
				ID:        shipID,
				Name:      warship.Name,
				Tier:      warship.Tier,
				Type:      model.NewShipType(warship.Type),
				Nation:    model.Nation(warship.Nation),
				IsPremium: warship.IsPremium,
			}
			mu.Unlock()
		}
		return pageTotal, nil
	}

	first := 1
	pageTotal, err := fetch(first)
	if err != nil {
		channel <- model.Result[model.Warships]{
			Error: err,
		}
		return
	}

	pages := makeRange(first+1, pageTotal+1)
	err = doParallel(pages, func(page int) error {
		_, err := fetch(page)
		return err
	})

	channel <- model.Result[model.Warships]{
		Value: warships,
		Error: err,
	}
}

func (f *WarshipFetcher) expectedStats(channel chan model.Result[response.ExpectedStats]) {
	es, err := f.numbers.ExpectedStats()

	channel <- model.Result[response.ExpectedStats]{
		Value: es,
		Error: err,
	}
}

func (f *WarshipFetcher) readCache() (warshipsCache, error) {
	cache, err := read[warshipsCache](f.db, f.localDataKeyName)
	return cache, failure.Translate(err, apperr.FetchShipError)
}

func (f *WarshipFetcher) saveCache(warships model.Warships, gameVersion string) error {
	cache := warshipsCache{
		warships:    warships,
		gameVersion: gameVersion,
	}
	return write(f.db, f.localDataKeyName, cache)
}

func (f *WarshipFetcher) toError(cache warshipsCache, errCache error, err error) (model.Warships, error) {
	if errCache != nil {
		err = errors.Join(errCache, failure.Translate(err, apperr.FetchShipError))
		return cache.warships, err
	}

	return cache.warships, nil
}
