//nolint:dupl
package service

import (
	"context"
	"sync"
	"wfs/backend/adapter"
	"wfs/backend/core"

	"github.com/samber/do/v2"
	"golang.org/x/sync/errgroup"
)

type StatsFetcher struct {
	wargamingClient adapter.WargamingClient
}

func NewStatsFetcher(i do.Injector) (*StatsFetcher, error) {
	return &StatsFetcher{
		wargamingClient: do.MustInvoke[adapter.WargamingClient](i),
	}, nil
}

func (f *StatsFetcher) FetchAll(
	ctx context.Context,
	accountIDs []core.AccountID,
) (core.AllPlayerShipStats, error) {
	result := make(core.AllPlayerShipStats)
	eg, egCtx := errgroup.WithContext(ctx)

	var mu sync.Mutex
	for _, accountID := range accountIDs {
		eg.Go(func() error {
			resp, err := f.wargamingClient.ShipsStats(egCtx, accountID)
			if err != nil {
				return err
			}

			mu.Lock()
			shipMap := make(core.PlayerShipStats)
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
