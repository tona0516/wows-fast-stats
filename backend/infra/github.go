package infra

import (
	"wfs/backend/vo"

	"github.com/cenkalti/backoff/v4"
)

type Github struct {
	latestReleaseClient APIClientInterface[vo.GHLatestRelease]
}

func NewGithub(config vo.GHConfig) *Github {
	return &Github{
		latestReleaseClient: NewAPIClient[vo.GHLatestRelease](
			config.BaseURL + "/repos/tona0516/wows-fast-stats/releases/latest",
		),
	}
}

func (g *Github) LatestRelease() (vo.GHLatestRelease, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)
	operation := func() (APIResponse[vo.GHLatestRelease], error) {
		return g.latestReleaseClient.GetRequest(make(map[string]string))
	}

	res, err := backoff.RetryWithData(operation, b)

	return res.Body, err
}
