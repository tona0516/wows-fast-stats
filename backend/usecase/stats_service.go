//nolint:dupl
package usecase

import (
	"sync"
	"wfs/backend/adapter"
	"wfs/backend/data"

	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
	"golang.org/x/sync/errgroup"
)

type statsService struct {
	wargamingClient adapter.WargamingClient
}

func NewStatsService(i do.Injector) (*statsService, error) {
	return &statsService{
		wargamingClient: do.MustInvoke[adapter.WargamingClient](i),
	}, nil
}

func (s *statsService) fetchAll(accountIDs []data.AccountID) (data.AllPlayerShipStats, error) {
	result := make(data.AllPlayerShipStats)
	var mu sync.Mutex

	eg := errgroup.Group{}

	for _, accountID := range accountIDs {
		eg.Go(func() error {
			resp, err := s.wargamingClient.ShipsStats(accountID)
			if err != nil {
				return failure.Wrap(err)
			}

			mu.Lock()
			shipMap := make(data.PlayerShipStats)
			for _, ship := range resp.Data[accountID] {
				shipMap[ship.ShipID] = ship
			}
			result[accountID] = shipMap
			mu.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, failure.Wrap(err)
	}

	return result, nil
}
