package usecase

import "wfs/backend/core"

type PollingResult struct {
	TempArenaInfo *core.TempArenaInfo
	Error         error
}
