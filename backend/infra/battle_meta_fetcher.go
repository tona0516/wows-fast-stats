package infra

import (
	"errors"
	"wfs/backend/apperr"
	"wfs/backend/domain/model"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
)

type BattleMetaFetcher struct {
	wargamingClient *req.Client
	cache           *model.BattleMeta
}

func NewBattleMetaFetcher(
	wargamingClient *req.Client,
) *BattleMetaFetcher {
	return &BattleMetaFetcher{
		wargamingClient: wargamingClient,
	}
}

func (f *BattleMetaFetcher) Fetch() (model.BattleMeta, error) {
	if f.cache != nil {
		return *f.cache, nil
	}

	arenaResultChan := make(chan model.Result[WGBattleArenas])
	typeResultChan := make(chan model.Result[WGBattleTypes])

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

func (f *BattleMetaFetcher) fetchBattleArenas(channel chan model.Result[WGBattleArenas]) {
	result := model.Result[WGBattleArenas]{}
	defer func() {
		channel <- result
	}()

	var body WGBattleArenasResponse
	_, err := f.wargamingClient.R().
		SetSuccessResult(&body).
		AddQueryParam("fields", WGBattleArenasResponse{}.Field()).
		AddQueryParam("language", "ja").
		Get("/wows/encyclopedia/battlearenas/")
	if err != nil {
		result.Error = failure.Wrap(err)
		return
	}

	result.Value = body.Data
}

func (f *BattleMetaFetcher) fetchBattleTypes(channel chan model.Result[WGBattleTypes]) {
	result := model.Result[WGBattleTypes]{}
	defer func() {
		channel <- result
	}()

	var body WGBattleTypesResponse
	_, err := f.wargamingClient.R().
		SetSuccessResult(&body).
		AddQueryParam("fields", WGBattleTypesResponse{}.Field()).
		AddQueryParam("language", "ja").
		Get("/wows/encyclopedia/battletypes/")
	if err != nil {
		result.Error = failure.Wrap(err)
		return
	}

	result.Value = body.Data
}
