package infra

import (
	"net/http"
	"time"
	"wfs/backend/apperr"
	"wfs/backend/data"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
)

type Github struct {
	url        string
	maxRetry   int
	timeoutSec int
}

func NewGithub(
	url string,
	maxRetry int,
	timeoutSec int,
) *Github {
	return &Github{
		url:        url,
		maxRetry:   maxRetry,
		timeoutSec: timeoutSec,
	}
}

func (g *Github) LatestRelease() (data.GHLatestRelease, error) {
	result := data.GHLatestRelease{}

	c := req.C().
		SetBaseURL(g.url).
		SetCommonRetryCount(g.maxRetry).
		SetTimeout(time.Duration(g.timeoutSec) * time.Second)

	resp, err := c.R().
		SetSuccessResult(&result).
		Get("/repos/tona0516/wows-fast-stats/releases/latest")

	if err != nil {
		return result, failure.Translate(err, apperr.GithubAPICheckUpdateError)
	}

	if resp.StatusCode != http.StatusOK {
		return result, failure.New(apperr.GithubAPICheckUpdateError, failure.Context{
			"status": resp.Status,
			"body":   resp.String(),
		})
	}

	return result, nil
}
