//nolint:dupl
package infra

import (
	"net/http"
	"strconv"
	"wfs/backend/apperr"
	"wfs/backend/domain"

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
	operation := func() (APIResponse[domain.GHLatestRelease], error) {
		return getRequest[domain.GHLatestRelease](
			g.config.URL + "/repos/tona0516/wows-fast-stats/releases/latest",
		)
	}

	res, err := backoff.RetryWithData(operation, b)
	if err != nil {
		return res.Body, failure.Wrap(err)
	}

	if res.StatusCode != http.StatusOK {
		return res.Body, failure.New(
			apperr.GithubAPICheckUpdateError,
			failure.Context{
				"status_code": strconv.Itoa(res.StatusCode),
				"body":        string(res.ByteBody),
			},
		)
	}

	return res.Body, nil
}
