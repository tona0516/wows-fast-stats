package infra

import (
	"net/http"
	"strconv"
	"wfs/backend/apperr"
	"wfs/backend/data"

	"github.com/go-resty/resty/v2"
	"github.com/morikuni/failure"
)

type Github struct {
	config apiConfig
}

func NewGithub(config apiConfig) *Github {
	return &Github{config: config}
}

func (g *Github) LatestRelease() (data.GHLatestRelease, error) {
	client := resty.New().
		SetTimeout(g.config.timeout).
		SetRetryCount(g.config.retryCount)

	var result data.GHLatestRelease
	resp, err := client.R().
		SetResult(&result).
		Get(g.config.baseURL + "/repos/tona0516/wows-fast-stats/releases/latest")

	errCtx := failure.Context{
		"url":         resp.Request.URL,
		"status_code": strconv.Itoa(resp.StatusCode()),
		"body":        string(resp.Body()),
	}

	if err != nil {
		return result, failure.New(apperr.GithubAPICheckUpdateError, failure.Messagef(err.Error()), errCtx)
	}

	if resp.StatusCode() != http.StatusOK {
		return result, failure.New(apperr.GithubAPICheckUpdateError, errCtx)
	}

	return result, nil
}
