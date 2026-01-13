package infra

import (
	"context"
	"wfs/backend/config"
	"wfs/backend/data"

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

func (c *GithubClient) LatestRelease(ctx context.Context) (data.GHLatestRelease, error) {
	result := data.GHLatestRelease{}
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get("/repos/tona0516/wows-fast-stats/releases/latest")
	if err != nil {
		return result, failure.Wrap(err)
	}
	if resp.IsErrorState() {
		return result, failure.Wrap(ErrErrorResponse)
	}

	return result, nil
}
