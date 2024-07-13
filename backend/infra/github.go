package infra

import (
	"net/http"
	"strconv"
	"time"
	"wfs/backend/apperr"
	"wfs/backend/data"

	"github.com/go-resty/resty/v2"
	"github.com/morikuni/failure"
)

type Github struct {
	baseURL string
}

func NewGithub(baseURL string) *Github {
	return &Github{baseURL: baseURL}
}

func (g *Github) LatestRelease() (data.GHLatestRelease, error) {
	client := resty.New().
		SetTimeout(5 * time.Second).
		SetRetryCount(2)

	var result data.GHLatestRelease
	resp, err := client.R().
		SetResult(&result).
		Get(g.baseURL + "/repos/tona0516/wows-fast-stats/releases/latest")

	errCtx := failure.Context{
		"url":         resp.Request.URL,
		"status_code": strconv.Itoa(resp.StatusCode()),
		"body":        string(resp.Body()),
	}
	if err != nil {
		return result, failure.Wrap(err, errCtx)
	}

	if resp.StatusCode() != http.StatusOK {
		return result, failure.New(apperr.GithubAPICheckUpdateError, errCtx)
	}

	return result, failure.Wrap(err)
}
