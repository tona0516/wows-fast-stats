//nolint:dupl
package usecase

import (
	"context"
	"sync"
	"wfs/backend/adapter"
	"wfs/backend/core"

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

func (s *badgeService) fetchAll(ctx context.Context, accountIDs []core.AccountID) (core.AllPlayerShipBadges, error) {
	result := make(core.AllPlayerShipBadges)
	eg, egCtx := errgroup.WithContext(ctx)

	var mu sync.Mutex
	for _, accountID := range accountIDs {
		eg.Go(func() error {
			resp, err := s.wargamingClient.ShipsBadges(egCtx, accountID)
			if err != nil {
				return err
			}

			mu.Lock()
			shipMap := make(core.PlayerShipBadges)
			for _, ship := range resp.Data[accountID] {
				shipMap[ship.ShipID] = ship
			}
			result[accountID] = shipMap
			mu.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}
