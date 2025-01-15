package infra

import (
	"errors"
	"sync"
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/domain/model"
	"wfs/backend/infra/response"

	"github.com/morikuni/failure"
)

type (
	shipStatsMap         map[int]response.WGShipsStatsData
	shipStatsAccountsMap map[int]shipStatsMap
)

type RawStatFetcher struct {
	wargaming Wargaming
}

func NewRawStatFetcher(
	wargaming Wargaming,
) *RawStatFetcher {
	return &RawStatFetcher{
		wargaming: wargaming,
	}
}

func (f *RawStatFetcher) Fetch(accountIDs []int) (model.RawStats, error) {
	accountInfoChan := make(chan data.Result[response.WGAccountInfo])
	shipStatsChan := make(chan data.Result[shipStatsAccountsMap])

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

func (f *RawStatFetcher) accountInfo(accountIDs []int, channel chan data.Result[response.WGAccountInfo]) {
	accountInfo, err := f.wargaming.accountInfo(accountIDs)
	channel <- data.Result[response.WGAccountInfo]{Value: accountInfo, Error: err}
}

func (f *RawStatFetcher) shipStats(
	accountIDs []int,
	channel chan data.Result[shipStatsAccountsMap],
) {
	shipStatsAccounts := make(shipStatsAccountsMap)
	var mu sync.Mutex

	err := doParallel(accountIDs, func(accountID int) error {
		wgShipStats, err := f.wargaming.shipsStats(accountID)
		if err != nil {
			return err
		}

		shipStats := make(shipStatsMap)
		for _, v := range wgShipStats[accountID] {
			shipStats[v.ShipID] = v
		}

		mu.Lock()
		shipStatsAccounts[accountID] = shipStats
		mu.Unlock()

		return nil
	})

	channel <- data.Result[shipStatsAccountsMap]{Value: shipStatsAccounts, Error: err}
}
