package repository

import "wfs/backend/domain/model"

type GithubInterface interface {
	LatestRelease() (model.GHLatestRelease, error)
}
