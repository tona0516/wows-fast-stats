package infra

import (
	"errors"
	"wfs/backend/apperr"
	"wfs/backend/domain/model"
	"wfs/backend/infra/response"
	"wfs/backend/infra/webapi"

	"github.com/morikuni/failure"
)

type BattleMetaFetcher struct {
	wargaming webapi.Wargaming
	cache     *model.BattleMeta
}

func NewBattleMetaFetcher(
	wargaming webapi.Wargaming,
) *BattleMetaFetcher {
	return &BattleMetaFetcher{
		wargaming: wargaming,
	}
}

func (f *BattleMetaFetcher) Fetch() (model.BattleMeta, error) {
	if f.cache != nil {
		return *f.cache, nil
	}

	arenaResultChan := make(chan model.Result[response.WGBattleArenas])
	typeResultChan := make(chan model.Result[response.WGBattleTypes])

	go f.fetchBattleArenas(arenaResultChan)
	go f.fetchBattleTypes(typeResultChan)

	var err error
	arenaResult := <-arenaResultChan
	err = errors.Join(err, arenaResult.Error)
	typeResult := <-typeResultChan
	err = errors.Join(err, typeResult.Error)
	if err != nil {
		return model.BattleMeta{}, failure.Translate(err, apperr.FetchBattleMetaError)
	}

	arenas := map[int]string{}
	for mapID, info := range arenaResult.Value {
		arenas[mapID] = info.Name
	}

	types := map[string]string{}
	for _type, info := range typeResult.Value {
		types[_type] = info.Name
	}

	return *model.NewBattleMeta(arenas, types), nil
}

func (f *BattleMetaFetcher) fetchBattleArenas(channel chan model.Result[response.WGBattleArenas]) {
	battleArenas, err := f.wargaming.BattleArenas()
	channel <- model.Result[response.WGBattleArenas]{Value: battleArenas, Error: err}
}

func (f *BattleMetaFetcher) fetchBattleTypes(channel chan model.Result[response.WGBattleTypes]) {
	battleTypes, err := f.wargaming.BattleTypes()
	channel <- model.Result[response.WGBattleTypes]{Value: battleTypes, Error: err}
}
