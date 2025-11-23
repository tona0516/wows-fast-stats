package repository

import "wfs/backend/data"

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type LocalFileInterface interface {
	TempArenaInfo(installPath string) (data.TempArenaInfo, error)
}
