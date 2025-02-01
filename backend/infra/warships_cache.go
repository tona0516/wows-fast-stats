package infra

import "wfs/backend/domain/model"

type warshipsCache struct {
	warships    model.Warships
	gameVersion string
}
