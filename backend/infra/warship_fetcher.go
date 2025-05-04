package infra

import (
	"encoding/json"
	"errors"
	"strconv"
	"sync"
	"wfs/backend/apperr"
	"wfs/backend/domain/model"

	"github.com/dgraph-io/badger/v4"
	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
)

type WarshipFetcher struct {
	db               *badger.DB
	wargamingClient  req.Client
	numbersClient    req.Client
	localDataKeyName string
}

func NewWarshipFetcher(
	db *badger.DB,
	wargamingClient req.Client,
	numbersClient req.Client,
) *WarshipFetcher {
	return &WarshipFetcher{
		db:               db,
		wargamingClient:  wargamingClient,
		numbersClient:    numbersClient,
		localDataKeyName: "warships",
	}
}

func (f *WarshipFetcher) Fetch() (model.Warships, error) {
	cache, errCache := f.readCache()

	currentGameVersion, err := f.fetchGameVersion()
	if err != nil {
		return f.toError(cache, errCache, err)
	}

	if currentGameVersion == cache.gameVersion {
		return cache.warships, nil
	}

	encycShipsChan := make(chan model.Result[model.Warships])
	expectedStatsChan := make(chan model.Result[NumbersExpectedStats])

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
	for shipID, ship := range expectedStats.Value.Data {
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
		res, err := f.fetchEncycShips(page)
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
		return res.Meta.PageTotal, nil
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

func (f *WarshipFetcher) expectedStats(channel chan model.Result[NumbersExpectedStats]) {
	result := model.Result[NumbersExpectedStats]{}
	defer func() {
		channel <- result
	}()

	resp, err := f.numbersClient.R().Get("/personal/rating/expected/json/")
	if err != nil {
		result.Error = failure.Wrap(err)
		return
	}

	var stats NumbersExpectedStats
	if err := json.Unmarshal(resp.Bytes(), &stats); err != nil {
		result.Error = failure.Wrap(err)
		return
	}

	result.Value = stats
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

func (f *WarshipFetcher) fetchGameVersion() (string, error) {
	resp, err := f.wargamingClient.R().
		AddQueryParam("fields", "game_version").
		Get("/wows/encyclopedia/info/")
	if err != nil {
		return "", failure.Wrap(err)
	}

	var result WGEncycInfo
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return "", failure.Wrap(err)
	}

	return result.Data.GameVersion, nil
}

func (f *WarshipFetcher) fetchEncycShips(pageNo int) (WGEncycShips, error) {
	var result WGEncycShips

	resp, err := f.wargamingClient.R().
		AddQueryParam("fields", WGEncycShips{}.Field()).
		AddQueryParam("language", "ja").
		AddQueryParam("page_no", strconv.Itoa(pageNo)).
		Get("/wows/encyclopedia/ships/")
	if err != nil {
		return result, failure.Wrap(err)
	}

	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return result, failure.Wrap(err)
	}

	return result, nil
}
