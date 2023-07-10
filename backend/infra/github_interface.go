package infra

import "wfs/backend/vo"

type GithubInterface interface {
	LatestRelease() (vo.GHLatestRelease, error)
}
