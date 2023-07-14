package infra

import (
	"wfs/backend/vo"

	"github.com/cenkalti/backoff/v4"
)

type Github struct {
	config              vo.RequestConfig
	latestReleaseClient APIClientInterface[vo.GHLatestRelease]
}

func NewGithub(config vo.RequestConfig) *Github {
	return &Github{
		config: config,
		latestReleaseClient: NewAPIClient[vo.GHLatestRelease](
			config.URL+"/repos/tona0516/wows-fast-stats/releases/latest",
			config.Retry,
		),
	}
}

func (g *Github) LatestRelease() (vo.GHLatestRelease, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), g.config.Retry)
	operation := func() (APIResponse[vo.GHLatestRelease], error) {
		return g.latestReleaseClient.GetRequest(make(map[string]string))
	}

	res, err := backoff.RetryWithData(operation, b)

	return res.Body, err
}
