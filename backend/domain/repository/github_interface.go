package repository

import "wfs/backend/domain/model"

//go:generate mockgen -source=$GOFILE -destination ../mock_$GOPACKAGE/$GOFILE -package mock_$GOPACKAGE
type GithubInterface interface {
	LatestRelease() (model.GHLatestRelease, error)
}
