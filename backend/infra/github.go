package infra

import (
	"net/http"
	"wfs/backend/apperr"
	"wfs/backend/domain"
	"wfs/backend/infra/webapi"

	"github.com/cenkalti/backoff/v4"
	"github.com/morikuni/failure"
)

type Github struct {
	config RequestConfig
}

func NewGithub(config RequestConfig) *Github {
	return &Github{config: config}
}

func (g *Github) LatestRelease() (domain.GHLatestRelease, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), g.config.Retry)
	operation := func() (webapi.Response[domain.GHLatestRelease], error) {
		res, err := webapi.GetRequest[domain.GHLatestRelease](
			g.config.URL+"/repos/tona0516/wows-fast-stats/releases/latest",
			g.config.Timeout,
		)
		if err != nil {
			return res, err
		}

		if res.StatusCode != http.StatusOK {
			return res, failure.New(apperr.GithubAPICheckUpdateError)
		}

		return res, nil
	}

	res, err := backoff.RetryWithData(operation, b)

	return res.Body, failure.Wrap(err, apperr.ToRequestErrorContext(res))
}
