package adapter

import "wfs/backend/domain"

type GithubInterface interface {
	LatestRelease() (domain.GHLatestRelease, error)
}
