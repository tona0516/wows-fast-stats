package infra

import (
	"wfs/internal/data"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type GithubApiClient interface {
	LatestRelease() (data.GHLatestRelease, error)
}

type githubApiClient struct {
	client *req.Client
}

func NewGithubApiClient(ac apiConfig) GithubApiClient {
	return &githubApiClient{
		client: req.C().
			SetBaseURL(ac.url).
			SetCommonRetryCount(ac.retryCount).
			SetTimeout(ac.timeout),
	}
}

func (c *githubApiClient) LatestRelease() (data.GHLatestRelease, error) {
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
