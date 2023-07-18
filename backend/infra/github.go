package infra

import (
	"wfs/backend/domain"

	"github.com/cenkalti/backoff/v4"
)

type Github struct {
	config RequestConfig
}

func NewGithub(config RequestConfig) *Github {
	return &Github{config: config}
}

func (g *Github) LatestRelease() (domain.GHLatestRelease, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), g.config.Retry)
	operation := func() (APIResponse[domain.GHLatestRelease], error) {
		return getRequest[domain.GHLatestRelease](
			g.config.URL+"/repos/tona0516/wows-fast-stats/releases/latest",
			make(map[string]string),
			g.config.Retry,
		)
	}

	res, err := backoff.RetryWithData(operation, b)

	return res.Body, err
}
