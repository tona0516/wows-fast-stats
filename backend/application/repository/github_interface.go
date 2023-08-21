package repository

import "wfs/backend/domain"

type GithubInterface interface {
	LatestRelease() (domain.GHLatestRelease, error)
}
