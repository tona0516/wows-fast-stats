package infra

import (
	"net/http"
	"strconv"
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
	operation := func() (webapi.Response[any, data.GHLatestRelease], error) {
		res, err := webapi.GetRequest[data.GHLatestRelease](
			g.config.URL+"/repos/tona0516/wows-fast-stats/releases/latest",
			g.config.Timeout,
			nil,
			g.config.Transport,
		)
		errCtx := failure.Context{
			"url":         res.Request.URL,
			"status_code": strconv.Itoa(res.StatusCode),
			"body":        string(res.BodyByte),
		}
		if err != nil {
			return res, failure.Wrap(err, errCtx)
		}

		if res.StatusCode != http.StatusOK {
			return res, failure.New(apperr.GithubAPICheckUpdateError, errCtx)
		}

		return res, nil
	}

	res, err := backoff.RetryWithData(operation, b)

	return res.Body, failure.Wrap(err)
}
