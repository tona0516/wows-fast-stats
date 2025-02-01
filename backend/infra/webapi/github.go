package webapi

import (
	"encoding/json"
	"wfs/backend/apperr"
	"wfs/backend/infra/response"

	"github.com/cenkalti/backoff/v4"
	"github.com/morikuni/failure"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type Github interface {
	LatestRelease() (response.GHLatestRelease, error)
}

type github struct {
	config RequestConfig
}

func NewGithub(config RequestConfig) Github {
	return &github{config: config}
}

func (g *github) LatestRelease() (response.GHLatestRelease, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), g.config.Retry)
	operation := func() (response.GHLatestRelease, error) {
		var result response.GHLatestRelease

		_, body, err := NewClient(g.config.URL,
			WithPath("/repos/tona0516/wows-fast-stats/releases/latest"),
			WithTimeout(g.config.Timeout),
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
