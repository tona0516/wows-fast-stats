package infra

import (
	"wfs/internal/data"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
)

type GithubClient struct {
	client *req.Client
}

func NewGithubClient(ac apiConfig) *GithubClient {
	return &GithubClient{
		client: req.C().
			SetBaseURL(ac.url).
			SetCommonRetryCount(ac.retryCount).
			SetTimeout(ac.timeout),
	}
}

func (c *GithubClient) LatestRelease() (data.GHLatestRelease, error) {
	result := data.GHLatestRelease{}
	resp, err := c.client.R().
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
