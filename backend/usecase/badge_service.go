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

type badgeService struct {
	wargamingClient adapter.WargamingClient
}

func NewBadgeService(i do.Injector) (*badgeService, error) {
	return &badgeService{
		wargamingClient: do.MustInvoke[adapter.WargamingClient](i),
	}, nil
}

func (s *badgeService) fetchAll(accountIDs []data.AccountID) (data.AllPlayerShipBadges, error) {
	result := make(data.AllPlayerShipBadges)
	var mu sync.Mutex

	eg := errgroup.Group{}
	for _, accountID := range accountIDs {
		eg.Go(func() error {
			resp, err := s.wargamingClient.ShipsBadges(accountID)
			if err != nil {
				return failure.Wrap(err)
			}

			mu.Lock()
			shipMap := make(data.PlayerShipBadges)
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
