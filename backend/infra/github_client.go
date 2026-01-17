package infra

import (
	"context"
	"wfs/backend/config"
	"wfs/backend/core"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
)

type GithubClient struct {
	client *req.Client
}

func NewGithubClient(i do.Injector) (*GithubClient, error) {
	config := do.MustInvoke[config.Config](i)
	return &GithubClient{
		client: req.C().
			SetBaseURL(config.GithubClient.URL).
			SetCommonRetryCount(config.GithubClient.RetryCount).
			SetTimeout(config.GithubClient.Timeout),
	}, nil
}

func (c *GithubClient) LatestRelease(ctx context.Context) (core.GHLatestRelease, error) {
	result := core.GHLatestRelease{}
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get("/repos/tona0516/wows-fast-stats/releases/latest")
	if err != nil {
		return result, failure.Translate(err, core.ErrGithubAPI)
	}
	if resp.IsErrorState() {
		return result, failure.New(core.ErrGithubAPI, failure.Context{
			"status_code": resp.Status,
			"body":        resp.String(),
		})
	}

	return result, nil
}
