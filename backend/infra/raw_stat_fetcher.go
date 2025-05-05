package infra

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"wfs/backend/apperr"
	"wfs/backend/domain/model"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
)

type (
	shipStatsMap         map[int]WGShipsStatsData
	shipStatsAccountsMap map[int]shipStatsMap
)

type RawStatFetcher struct {
	wargamingClient *req.Client
}

func NewRawStatFetcher(
	wargamingClient *req.Client,
) *RawStatFetcher {
	return &RawStatFetcher{
		wargamingClient: wargamingClient,
	}
}

func (f *RawStatFetcher) Fetch(accountIDs []int) (model.RawStats, error) {
	accountInfoChan := make(chan model.Result[WGAccountInfo])
	shipStatsChan := make(chan model.Result[shipStatsAccountsMap])

	go f.accountInfo(accountIDs, accountInfoChan)
	go f.shipStats(accountIDs, shipStatsChan)

	var err error
	accountInfoResult := <-accountInfoChan
	err = errors.Join(err, accountInfoResult.Error)

	shipStatsResult := <-shipStatsChan
	err = errors.Join(err, shipStatsResult.Error)
	if err != nil {
		return nil, failure.Translate(err, apperr.FetchRawStatError)
	}

	rawStats := make(model.RawStats)
	for accountID, info := range accountInfoResult.Value {
		rawStat := model.RawStat{
			Ship: map[int]model.ShipStat{},
			Overall: model.OverallStat{
				IsHidden: info.HiddenProfile,
				Pvp:      model.OverallStatsValues(info.Statistics.Pvp),
				PvpDiv2:  model.OverallStatsValues(info.Statistics.PvpDiv2),
				PvpDiv3:  model.OverallStatsValues(info.Statistics.PvpDiv3),
				PvpSolo:  model.OverallStatsValues(info.Statistics.PvpSolo),
				RankSolo: model.OverallStatsValues(info.Statistics.RankSolo),
			},
		}

		shipsStats := shipStatsResult.Value[accountID]
		for _, v := range shipsStats {
			rawStat.Ship[v.ShipID] = model.ShipStat{
				Pvp:      model.ShipStatsValues(v.Pvp),
				PvpSolo:  model.ShipStatsValues(v.PvpSolo),
				PvpDiv2:  struct{ Battles uint }(v.PvpDiv2),
				PvpDiv3:  struct{ Battles uint }(v.PvpDiv3),
				RankSolo: model.ShipStatsValues(v.RankSolo),
			}
		}

		rawStats[accountID] = rawStat
	}

	return rawStats, nil
}

func (f *RawStatFetcher) accountInfo(accountIDs []int, channel chan model.Result[WGAccountInfo]) {
	strAccountIDs := make([]string, len(accountIDs))
	for i, v := range accountIDs {
		strAccountIDs[i] = strconv.Itoa(v)
	}

	var result model.Result[WGAccountInfo]
	defer func() {
		channel <- result
	}()

	var body WGAccountInfoResponse
	_, err := f.wargamingClient.R().
		SetSuccessResult(&body).
		AddQueryParam("account_id", strings.Join(strAccountIDs, ",")).
		AddQueryParam("fields", WGAccountInfoResponse{}.Field()).
		AddQueryParam("extra", strings.Join([]string{
			"statistics.pvp_solo",
			"statistics.pvp_div2",
			"statistics.pvp_div3",
			"statistics.rank_solo",
		}, ",")).
		Get("/wows/account/info/")
	if err != nil {
		result.Error = failure.Wrap(err)
		return
	}

	result.Value = body.Data
}

func (f *RawStatFetcher) shipStats(
	accountIDs []int,
	channel chan model.Result[shipStatsAccountsMap],
) {
	shipStatsAccounts := make(shipStatsAccountsMap)
	var mu sync.Mutex

	err := doParallel(accountIDs, func(accountID int) error {
		var body WGShipsStatsResponse
		_, err := f.wargamingClient.R().
			SetSuccessResult(&body).
			AddQueryParam("account_id", strconv.Itoa(accountID)).
			AddQueryParam("fields", WGShipsStatsResponse{}.Field()).
			AddQueryParam("extra", strings.Join([]string{
				"pvp_solo",
				"pvp_div2",
				"pvp_div3",
				"rank_solo",
			}, ",")).
			Get("/wows/ships/stats/")
		if err != nil {
			return failure.Wrap(err)
		}

		shipStats := make(shipStatsMap)
		for _, v := range body.Data[accountID] {
			shipStats[v.ShipID] = v
		}

		mu.Lock()
		shipStatsAccounts[accountID] = shipStats
		mu.Unlock()

		return nil
	})

	channel <- model.Result[shipStatsAccountsMap]{Value: shipStatsAccounts, Error: err}
}
