package repository

import "wfs/internal/data"

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type GithubInterface interface {
	LatestRelease() (data.GHLatestRelease, error)
}
