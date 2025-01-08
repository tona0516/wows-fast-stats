package infra

import (
	"encoding/json"
	"wfs/backend/apperr"
	"wfs/backend/data"
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

func (g *Github) LatestRelease() (data.GHLatestRelease, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), g.config.Retry)
	operation := func() (data.GHLatestRelease, error) {
		var result data.GHLatestRelease

		_, body, err := webapi.NewClient(g.config.URL,
			webapi.WithPath("/repos/tona0516/wows-fast-stats/releases/latest"),
			webapi.WithTimeout(g.config.Timeout),
		).GET()
		if err != nil {
			return result, err
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return result, failure.Wrap(err)
		}

		return result, nil
	}

	res, err := backoff.RetryWithData(operation, b)
	if err != nil {
		return res, failure.Translate(err, apperr.GithubAPICheckUpdateError)
	}

	return res, nil
}
